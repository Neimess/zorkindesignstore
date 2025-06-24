package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/Neimess/zorkin-store-project/pkg/database"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
	"log/slog"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrCategoryNotFound = errors.New("category not found")
)

// ProductRepository handles CRUD for products and their attributes.
type ProductRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

// NewProductRepository creates a new repo instance or panics if db is nil.
func NewProductRepository(db *sqlx.DB, log *slog.Logger) *ProductRepository {
	if db == nil {
		panic("NewProductRepository: db is nil")
	}
	return &ProductRepository{
		db:  db,
		log: logger.WithComponent(log, "repo.product"),
	}
}

// Create inserts a product and populates its ID and CreatedAt.
func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) (int64, error) {
	const insertSQL = `
        INSERT INTO products(name, price, description, category_id, image_url)
        VALUES ($1,$2,$3,$4,$5)
        RETURNING product_id, created_at`

	var (
		id      int64
		created time.Time
	)

	err := r.withQuery(ctx, insertSQL, func() error {
		return r.db.QueryRowContext(ctx, insertSQL,
			p.Name, p.Price, p.Description, p.CategoryID, p.ImageURL,
		).Scan(&id, &created)
	})
	if err != nil {
		return 0, r.wrapError(err, "create product")
	}

	p.ID, p.CreatedAt = id, created
	return id, nil
}

// CreateWithAttrs inserts a product and its attributes in a transaction.
func (r *ProductRepository) CreateWithAttrs(ctx context.Context, p *domain.Product) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer r.rollbackOnPanicOrError(tx, &err)

	prodID, err := r.insertProductTx(ctx, tx, p)
	if err != nil {
		return 0, err
	}

	if len(p.Attributes) > 0 {
		if err := r.insertAttributesTx(ctx, tx, prodID, p.Attributes); err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}

	p.ID = prodID
	return prodID, nil
}

// GetWithAttrs fetches product and its attributes; returns ErrProductNotFound if missing.
func (r *ProductRepository) GetWithAttrs(ctx context.Context, id int64) (*domain.Product, error) {
	prod, err := r.fetchProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	attrs, err := r.fetchAttributes(ctx, id)
	if err != nil {
		return nil, err
	}
	prod.Attributes = attrs
	return prod, nil
}

// ListByCategory returns products in a category with their attributes.
func (r *ProductRepository) ListByCategory(ctx context.Context, catID int64) ([]domain.Product, error) {
	const listSQL = `
        SELECT product_id,name,price,description,category_id,image_url,created_at
        FROM products WHERE category_id=$1`

	var products []domain.Product
	if err := r.withQuery(ctx, listSQL, func() error {
		return r.db.SelectContext(ctx, &products, listSQL, catID)
	}); err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	if len(products) == 0 {
		return products, nil
	}

	ids := make([]int64, len(products))
	for i, p := range products {
		ids[i] = p.ID
	}

	attrRows, err := r.fetchAttributesBatch(ctx, ids)
	if err != nil {
		return nil, err
	}

	// Map attributes to products
	byID := make(map[int64]*domain.Product, len(products))
	for i := range products {
		byID[products[i].ID] = &products[i]
	}
	for _, row := range attrRows {
		attr := domain.ProductAttribute{
			AttributeID: row.AttributeID,
			Value:       row.Value,
			Attribute: domain.Attribute{
				ID:           row.AttributeID,
				Name:         row.Name,
				Slug:         row.Slug,
				Unit:         optionalString(row.Unit),
				IsFilterable: row.IsFilterable,
			},
		}
		byID[row.ProductID].Attributes = append(byID[row.ProductID].Attributes, attr)
	}
	return products, nil
}

//----------------------
// Internal Helpers
//----------------------

func (r *ProductRepository) withQuery(ctx context.Context, sql string, fn func() error) error {
	err := database.WithQuery(ctx, r.log, sql, fn)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return err
}

func (r *ProductRepository) wrapError(err error, op string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolationCode {
		return ErrInvalidForeignKey
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrProductNotFound
	}
	return fmt.Errorf("%s: %w", op, err)
}

func (r *ProductRepository) rollbackOnPanicOrError(tx *sqlx.Tx, errPtr *error) {
	if rec := recover(); rec != nil {
		_ = tx.Rollback()
		panic(rec)
	} else if *errPtr != nil {
		_ = tx.Rollback()
	}
}

func (r *ProductRepository) insertProductTx(ctx context.Context, tx *sqlx.Tx, p *domain.Product) (int64, error) {
	const sql = `
        INSERT INTO products(name, price, description, category_id, image_url)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING product_id`

	var id int64
	if err := database.WithQuery(ctx, r.log, sql, func() error {
		return tx.QueryRowContext(ctx, sql,
			p.Name, p.Price, p.Description, p.CategoryID, p.ImageURL,
		).Scan(&id)
	}); err != nil {
		return 0, r.wrapError(err, "insert product")
	}
	return id, nil
}

func (r *ProductRepository) insertAttributesTx(ctx context.Context, tx *sqlx.Tx, prodID int64, attrs []domain.ProductAttribute) error {
	builder := sq.Insert("product_attributes").Columns("product_id", "attribute_id", "value").PlaceholderFormat(sq.Dollar)
	for _, a := range attrs {
		builder = builder.Values(prodID, a.AttributeID, a.Value)
	}
	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build attrs query: %w", err)
	}
	if err := database.WithQuery(ctx, r.log, sqlStr, func() error {
		_, err := tx.ExecContext(ctx, sqlStr, args...)
		return err
	}); err != nil {
		return r.wrapError(err, "insert attributes")
	}
	return nil
}

func (r *ProductRepository) fetchProduct(ctx context.Context, id int64) (*domain.Product, error) {
	const sql = `
        SELECT product_id, name, price, description, category_id, image_url, created_at
        FROM products WHERE product_id = $1`

	var prod domain.Product
	if err := database.WithQuery(ctx, r.log, sql, func() error {
		return r.db.GetContext(ctx, &prod, sql, id)
	}); err != nil {
		return nil, r.wrapError(err, "get product")
	}
	return &prod, nil
}

type attrRow struct {
	ProductID    int64          `db:"product_id"`
	AttributeID  int64          `db:"attribute_id"`
	Value        string         `db:"value"`
	Name         string         `db:"name"`
	Slug         string         `db:"slug"`
	Unit         sql.NullString `db:"unit"`
	IsFilterable bool           `db:"is_filterable"`
}

func (r *ProductRepository) fetchAttributes(ctx context.Context, prodID int64) ([]domain.ProductAttribute, error) {
	const sql = `
        SELECT pa.attribute_id, pa.value, a.name, a.slug, a.unit, a.is_filterable
        FROM product_attributes pa
        JOIN attributes a ON a.attribute_id = pa.attribute_id
        WHERE pa.product_id = $1`

	rows, err := r.db.QueryxContext(ctx, sql, prodID)
	if err != nil {
		return nil, fmt.Errorf("query attrs: %w", err)
	}
	defer rows.Close()

	var result []domain.ProductAttribute
	for rows.Next() {
		var ar attrRow
		if err := rows.StructScan(&ar); err != nil {
			return nil, fmt.Errorf("scan attribute: %w", err)
		}
		attr := domain.ProductAttribute{
			AttributeID: ar.AttributeID,
			Value:       ar.Value,
			Attribute: domain.Attribute{
				ID:           ar.AttributeID,
				Name:         ar.Name,
				Slug:         ar.Slug,
				Unit:         optionalString(ar.Unit),
				IsFilterable: ar.IsFilterable,
			},
		}
		result = append(result, attr)
	}
	return result, rows.Err()
}

func (r *ProductRepository) fetchAttributesBatch(ctx context.Context, ids []int64) ([]attrRow, error) {
	const sql = `
        SELECT pa.product_id, pa.attribute_id, pa.value,
               a.name, a.slug, a.unit, a.is_filterable
        FROM product_attributes pa
        JOIN attributes a USING(attribute_id)
        WHERE pa.product_id = ANY($1)`

	var rows []attrRow
	if err := database.WithQuery(ctx, r.log, sql, func() error {
		return r.db.SelectContext(ctx, &rows, sql, pq.Array(ids))
	}); err != nil {
		return nil, fmt.Errorf("attrs by category: %w", err)
	}
	return rows, nil
}

func optionalString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
