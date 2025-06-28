package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	"github.com/Neimess/zorkin-store-project/internal/domain/product"
	repoError "github.com/Neimess/zorkin-store-project/internal/infrastructure/error"
	"github.com/Neimess/zorkin-store-project/pkg/database"
	tx "github.com/Neimess/zorkin-store-project/pkg/database/tx"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrCategoryNotFound = errors.New("category not found")
)

// PGProductRepository handles CRUD for products and their attributes.
type PGProductRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

// NewProductRepository creates a new repo instance or panics if db is nil.
func NewPGProductRepository(db *sqlx.DB, log *slog.Logger) *PGProductRepository {
	if db == nil {
		panic("NewProductRepository: db is nil")
	}
	return &PGProductRepository{
		db:  db,
		log: log.With("component", "repo.product"),
	}
}

// Create inserts a product and populates its ID and CreatedAt.
func (r *PGProductRepository) Create(ctx context.Context, p *product.Product) (int64, error) {
	const insertSQL = `
        INSERT INTO products(name, price, description, category_id, image_url)
        VALUES ($1,$2,$3,$4,$5)
        RETURNING product_id, created_at`

	var (
		id      int64
		created time.Time
	)

	err := database.WithQuery(ctx, slog.Default().With(), insertSQL, func() error {
		return r.db.QueryRowContext(ctx, insertSQL,
			p.Name, p.Price, p.Description, p.CategoryID, p.ImageURL,
		).Scan(&id, &created)
	})
	if err != nil {
		return 0, r.mapPostgreSQLError(err)
	}

	p.ID, p.CreatedAt = id, created
	return id, nil
}

// CreateWithAttrs inserts a product and its attributes in a transaction.
func (r *PGProductRepository) CreateWithAttrs(ctx context.Context, p *product.Product) (int64, error) {
	return tx.RunInTx(ctx, r.db, func(tx *sqlx.Tx) (int64, error) {
		// insert main product
		prodID, err := r.insertProductTx(ctx, tx, p)
		if err != nil {
			return 0, err
		}

		// insert attributes if any
		if len(p.Attributes) > 0 {
			if err := r.insertAttributesTx(ctx, tx, prodID, p.Attributes); err != nil {
				return 0, err
			}
		}

		p.ID = prodID
		return prodID, nil
	})
}

// GetWithAttrs fetches product and its attributes; returns ErrProductNotFound if missing.
func (r *PGProductRepository) GetWithAttrs(ctx context.Context, id int64) (*product.Product, error) {
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
func (r *PGProductRepository) ListByCategory(ctx context.Context, catID int64) ([]product.Product, error) {
	const listSQL = `
        SELECT product_id,name,price,description,category_id,image_url,created_at
        FROM products WHERE category_id=$1`

	var products []product.Product
	err := r.withQuery(ctx, listSQL, func() error {
		return r.db.SelectContext(ctx, &products, listSQL, catID)
	})
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
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

	byID := make(map[int64]*product.Product, len(products))
	for i := range products {
		byID[products[i].ID] = &products[i]
	}
	for _, row := range attrRows {
		attr := product.ProductAttribute{
			AttributeID: row.AttributeID,
			Value:       row.Value,
			Attribute: attr.Attribute{
				ID:   row.AttributeID,
				Name: row.Name,
				Unit: optionalString(row.Unit),
			},
		}
		byID[row.ProductID].Attributes = append(byID[row.ProductID].Attributes, attr)
	}
	return products, nil
}

//----------------------
// Internal Helpers
//----------------------

func (r *PGProductRepository) insertProductTx(ctx context.Context, tx *sqlx.Tx, p *product.Product) (int64, error) {
	const sqlStr = `
        INSERT INTO products(name, price, description, category_id, image_url)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING product_id`

	var id int64
	if err := r.withQuery(ctx, sqlStr, func() error {
		return tx.QueryRowContext(ctx, sqlStr,
			p.Name, p.Price, p.Description, p.CategoryID, p.ImageURL,
		).Scan(&id)
	}); err != nil {
		return 0, r.mapPostgreSQLError(err)
	}
	return id, nil
}

func (r *PGProductRepository) insertAttributesTx(ctx context.Context, tx *sqlx.Tx, prodID int64, attrs []product.ProductAttribute) error {
	builder := sq.Insert("product_attributes").Columns("product_id", "attribute_id", "value").PlaceholderFormat(sq.Dollar)
	for _, a := range attrs {
		builder = builder.Values(prodID, a.AttributeID, a.Value)
	}
	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build attrs query: %w", err)
	}
	if err := r.withQuery(ctx, sqlStr, func() error {
		_, execErr := tx.ExecContext(ctx, sqlStr, args...)
		return execErr
	}); err != nil {
		return r.mapPostgreSQLError(err)
	}
	return nil
}

type attrRow struct {
	ProductID   int64          `db:"product_id"`
	AttributeID int64          `db:"attribute_id"`
	Value       string         `db:"value"`
	Name        string         `db:"name"`
	Unit        sql.NullString `db:"unit"`
}

func (r *PGProductRepository) fetchProduct(ctx context.Context, id int64) (*product.Product, error) {
	const sqlStr = `
        SELECT product_id, name, price, description, category_id, image_url, created_at
        FROM products WHERE product_id = $1`

	var prod product.Product
	err := r.withQuery(ctx, sqlStr, func() error {
		return r.db.GetContext(ctx, &prod, sqlStr, id)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, r.mapPostgreSQLError(err)
	}
	return &prod, nil
}

func (r *PGProductRepository) fetchAttributes(ctx context.Context, prodID int64) ([]product.ProductAttribute, error) {
	const sqlStr = `
        SELECT pa.attribute_id, pa.value, a.name, a.unit
        FROM product_attributes pa
        JOIN attributes a ON a.attribute_id = pa.attribute_id
        WHERE pa.product_id = $1`

	rows, err := r.db.QueryxContext(ctx, sqlStr, prodID)
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			slog.Error("rows.Close failed", slog.Any("error", cerr))
		}
	}()

	var result []product.ProductAttribute
	for rows.Next() {
		var ar attrRow
		if err := rows.StructScan(&ar); err != nil {
			return nil, fmt.Errorf("scan attribute: %w", err)
		}
		attr := product.ProductAttribute{
			AttributeID: ar.AttributeID,
			Value:       ar.Value,
			Attribute: attr.Attribute{
				ID:   ar.AttributeID,
				Name: ar.Name,
				Unit: optionalString(ar.Unit),
			},
		}
		result = append(result, attr)
	}
	return result, rows.Err()
}

func (r *PGProductRepository) fetchAttributesBatch(ctx context.Context, ids []int64) ([]attrRow, error) {
	const sqlStr = `
        SELECT pa.product_id, pa.attribute_id, pa.value,
               a.name, a.unit
        FROM product_attributes pa
        JOIN attributes a USING(attribute_id)
        WHERE pa.product_id = ANY($1)`

	var rows []attrRow
	if err := r.withQuery(ctx, sqlStr, func() error {
		return r.db.SelectContext(ctx, &rows, sqlStr, pq.Array(ids))
	}); err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	return rows, nil
}

func optionalString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func (r *PGProductRepository) withQuery(ctx context.Context, query string, fn func() error, extras ...slog.Attr) error {
	return database.WithQuery(ctx, r.log, query, fn, extras...)
}

func (r *PGProductRepository) mapPostgreSQLError(err error) error {
	return repoError.MapPostgreSQLError(r.log, err)
}
