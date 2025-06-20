package repository

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/jmoiron/sqlx"
)

//go:embed sql/*.sql
var queries embed.FS



type ProductRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewProductRepository(db *sqlx.DB, logger *slog.Logger) *ProductRepository {
	return &ProductRepository{db: db, log: logger.With("component", "repo.product")}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p *domain.Product) (int64, error) {
	const query = `
	INSERT INTO products(name, price, description, category_id, image_url)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING product_id, created_at
	`
	done := r.dbLog(ctx, query)

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
		return 0, fmt.Errorf("create product: %w", err)
	}
	p.ID = id
	p.CreatedAt = createdAt
	return id, nil

}

// GetProduct fetches a product by its id along with its attributes.
// Returns ErrProductNotFound if the row does not exist.
func (r *ProductRepository) GetProductWithAttrs(ctx context.Context, id int64) (*domain.Product, error) {
	const query = `
        SELECT product_id, name, price, description, category_id, image_url, created_at
        FROM products
        WHERE product_id = $1
    `
	done := r.dbLog(ctx, query)
	var p domain.Product
	err := r.db.GetContext(ctx, &p, query, id)
	done(err, slog.Int64("product_id", id))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return nil, err
		default:
			return nil, fmt.Errorf("get product failed: id=%d, category=%d, name=%q: %w",
				id, p.CategoryID, p.Name, err,
			)
		}
	}

	attrSQL, err := queries.ReadFile("sql/productAttributesQuery.sql")
	if err != nil {
		return nil, fmt.Errorf("read productAttributesQuery.sql: %w", err)
	}

	attrDone := r.dbLog(ctx, string(attrSQL))
	var attrs []domain.ProductAttribute
	err = r.db.SelectContext(ctx, &attrs, string(attrSQL), p.ID)
	attrDone(err)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, fmt.Errorf("fetch product attributes for product_id=%d: %w", p.ID, err)
	}
	p.Attributes = attrs
	return &p, nil
}
