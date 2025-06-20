package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Neimess/zorkin-store-project/internal/domain"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewProductRepository(db *sqlx.DB, log *slog.Logger) *ProductRepository {
	if db == nil {
		panic("NewProductRepository: db is nil")
	}
	return &ProductRepository{
		db:  db,
		log: logger.WithComponent(log, "repo.category"),
	}
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) (int64, error) {
	const query = `
	INSERT INTO products(name, price, description, category_id, image_url)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING product_id, created_at
	`
	done := dbLog(ctx, r.log, query)

	var id int64
	var createdAt time.Time

	err := r.db.QueryRowContext(ctx, query,
		p.Name, p.Price, p.Description, p.CategoryID, p.ImageURL,
	).Scan(&id, &createdAt)

	done(err, slog.Int64("product_id", id))
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return 0, err
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolationCode {
			return 0, ErrInvalidForeignKey
		}
		return 0, fmt.Errorf("create product: %w", err)
	}
	p.ID = id
	p.CreatedAt = createdAt
	return id, nil

}

func (r *ProductRepository) CreateWithAttrs(ctx context.Context, p *domain.Product) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	const insertProduct = `
        INSERT INTO products (name, price, description, category_id, image_url)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING product_id
    `
	var productID int64
	if err := tx.QueryRowContext(ctx, insertProduct,
		p.Name, p.Price, p.Description, p.CategoryID, p.ImageURL,
	).Scan(&productID); err != nil {
		return 0, fmt.Errorf("insert product: %w", err)
	}

	if len(p.Attributes) > 0 {
		builder := sq.
			Insert("product_attributes").
			Columns("product_id", "attribute_id", "value").
			PlaceholderFormat(sq.Dollar)

		for _, attr := range p.Attributes {
			builder = builder.Values(
				productID,
				attr.AttributeID,
				attr.Value,
			)
		}

		sqlAttrs, args, err := builder.ToSql()
		if err != nil {
			return 0, fmt.Errorf("build attrs query: %w", err)
		}
		if _, err := tx.ExecContext(ctx, sqlAttrs, args...); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolationCode {
				return 0, ErrInvalidForeignKey
			}
			return 0, fmt.Errorf("insert attributes: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}
	p.ID = productID
	return productID, nil
}

// GetProduct fetches a product by its id along with its attributes.
// Returns ErrProductNotFound if the row does not exist.
func (r *ProductRepository) GetWithAttrs(ctx context.Context, id int64) (*domain.Product, error) {
	const query = `
        SELECT product_id, name, price, description, category_id, image_url, created_at
        FROM products
        WHERE product_id = $1
    `

	done := dbLog(ctx, r.log, query)

	var p domain.Product
	err := r.db.GetContext(ctx, &p, query, id)
	done(err, slog.Int64("product_id", id))

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrProductNotFound
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return nil, err
		default:
			return nil, fmt.Errorf("get product failed: id=%d, category=%d, name=%q: %w",
				id, p.CategoryID, p.Name, err,
			)
		}
	}

	const q2 = `
        SELECT
          pa.attribute_id,
          pa.value,
          a.name, a.slug, a.data_type, a.unit, a.is_filterable
        FROM product_attributes pa
        JOIN attributes a ON a.attribute_id = pa.attribute_id
        WHERE pa.product_id = $1
    `
	rows, err := r.db.QueryxContext(ctx, q2, id)
	if err != nil {
		return nil, fmt.Errorf("query attributes: %w", err)
	}
	defer rows.Close()

	attrDone := dbLog(ctx, r.log, q2)
	var attrs []domain.ProductAttribute
	err = r.db.SelectContext(ctx, &attrs, q2, p.ID)
	attrDone(err)
	for rows.Next() {
		var pa domain.ProductAttribute
		var meta struct {
			Name         string
			Slug         string
			Unit         sql.NullString
			IsFilterable bool
		}
		if err := rows.Scan(
			&pa.AttributeID,
			&pa.Value,
			&meta.Name,
			&meta.Slug,

			&meta.Unit,
			&meta.IsFilterable,
		); err != nil {
			return nil, fmt.Errorf("scan attribute row: %w", err)
		}
		pa.Attribute = domain.Attribute{
			ID:           pa.AttributeID,
			Name:         meta.Name,
			Slug:         meta.Slug,
			Unit:         optionalString(meta.Unit),
			IsFilterable: meta.IsFilterable,
		}
		attrs = append(attrs, pa)
	}
	p.Attributes = attrs
	return &p, nil
}

func optionalString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
