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
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
)

//go:embed sql/*.sql
var queries embed.FS

var (
	ErrProductIsNil      = errors.New("product is nil")
	ErrProductNotFound   = errors.New("product not found")
	ErrCategoryNotExists = errors.New("referenced category does not exist")
)

type ProductRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewProductRepository(db *sqlx.DB, logger *slog.Logger) *ProductRepository {
	return &ProductRepository{db: db, log: logger.With("component", "repo.product")}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p *domain.Product) (int64, error) {
	const op = "repository.psql.product.create_product"
	log := r.log.With("op", op)
	if p == nil {
		return 0, ErrProductIsNil
	}

	const query = `
	INSERT INTO products(name, price, description, category_id, image_url)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING product_id, created_at
	`
	var (
		id         int64
		created_at time.Time
	)

	log.Debug("exec query", slog.String("query", query))
	err := r.db.QueryRowContext(
		ctx,
		query,
		p.Name,
		p.Price,
		p.Description,
		p.CategoryID,
		p.ImageURL,
	).Scan(&id, &created_at)

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return 0, err
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return 0, ErrCategoryNotExists
			case "23505":
				return 0, fmt.Errorf("product already exists: %w", err)
			}
		}

		return 0, fmt.Errorf("create product: %w", err)
	}

	p.ID = id
	p.CreatedAt = created_at
	return id, nil

}

// GetProduct fetches a product by its id along with its attributes.
// Returns ErrProductNotFound if the row does not exist.
func (r *ProductRepository) GetProductWithAttrs(ctx context.Context, id int64) (*domain.Product, error) {
	const op = "repository.psql.product.create_product"
	log := r.log.With("op", op)
	const query = `
        SELECT product_id, name, price, description, category_id, image_url, created_at
        FROM products
        WHERE product_id = $1
    `
	log.Debug("exec query", slog.String("query", query))
	var p domain.Product
	if err := r.db.GetContext(ctx, &p, query, id); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrProductNotFound
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return nil, err
		default:
			return nil, fmt.Errorf("get product by id=%d: %w", id, err)
		}
	}
	log.Debug(fmt.Sprintf("select product with id = %d", id), "product", p)
	attrQuery, err := queries.ReadFile("sql/productAttributesQuery.sql")
	if err != nil {
		return nil, fmt.Errorf("read productAttributesQuery.sql: %w", err)
	}

	log.Debug("exec query", slog.String("query", string(attrQuery)))
	var attrs []domain.ProductAttribute
	if err := r.db.SelectContext(ctx, &attrs, string(attrQuery), p.ID); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, fmt.Errorf("fetch product attributes for product_id=%d: %w", p.ID, err)
	}
	log.Debug("getted attrs", "attrs", attrs)
	p.Attributes = attrs
	log.Debug("final product with attrs", "product", p)
	return &p, nil
}
