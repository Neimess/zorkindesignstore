package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	repository "github.com/Neimess/zorkin-store-project/internal/repository/psql"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
)

var (
	ErrNoProductsFound = errors.New("products not found")
	ErrProductNotFound  = errors.New("product not found")
	ErrBadCategoryID    = errors.New("bad category ID")
	ErrInvalidAttribute = errors.New("invalid product attribute")
)

type ProductRepository interface {
	Create(ctx context.Context, p *domain.Product) (int64, error)
	CreateWithAttrs(ctx context.Context, p *domain.Product) (int64, error)
	GetWithAttrs(ctx context.Context, id int64) (*domain.Product, error)
	ListByCategory(ctx context.Context, catID int64) ([]domain.Product, error)
}

type ProductService struct {
	repo ProductRepository
	log  *slog.Logger
}

func NewProductService(repo ProductRepository, log *slog.Logger) *ProductService {
	if repo == nil {
		panic("product service: repo is nil")
	}
	return &ProductService{
		repo: repo,
		log:  logger.WithComponent(log, "service.product"),
	}
}

func (ps *ProductService) Create(ctx context.Context, product *domain.Product) (int64, error) {
	const op = "service.product.Create"
	log := ps.log.With("op", op)

	id, err := ps.repo.Create(ctx, product)

	switch {
	case errors.Is(err, repository.ErrInvalidForeignKey):
		log.Warn("invalid category", slog.Int64("category_id", product.CategoryID))
		return 0, ErrBadCategoryID
	case err != nil:
		log.Error("repo failed", slog.Int64("product_id", product.ID), slog.Any("error", err))
		return 0, fmt.Errorf("%s: %w", op, err)
	default:
		log.Info("product created",
			slog.Int64("product_id", id))
	}
	return id, nil
}

func (ps *ProductService) CreateWithAttrs(ctx context.Context, product *domain.Product) (int64, error) {
	const op = "service.product.CreateWithAttributes"
	log := ps.log.With("op", op)

	id, err := ps.repo.CreateWithAttrs(ctx, product)
	switch {
	case errors.Is(err, repository.ErrInvalidForeignKey):
		log.Warn("invalid category", slog.Int64("category_id", product.CategoryID))
		return 0, ErrBadCategoryID
	case err != nil:
		log.Error("repo failed", slog.Int64("product_id", product.ID), slog.Any("error", err))
		return 0, fmt.Errorf("%s: %w", op, err)
	default:
		log.Info("product with attributes created",
			slog.Int64("product_id", id))
	}
	return id, nil

}

// GetDetailed returns domain.Product with it's domain.ProductAttribute
func (ps *ProductService) GetDetailed(ctx context.Context, id int64) (*domain.Product, error) {
	const op = "service.product.GetDetailed"
	log := ps.log.With("op", op)
	product, err := ps.repo.GetWithAttrs(ctx, id)

	switch {
	case errors.Is(err, repository.ErrProductNotFound):
		return nil, ErrProductNotFound
	case err != nil:
		log.Error("repo failed", slog.Int64("product_id", id), slog.Any("err", err))
		return nil, err
	}
	log.Info("product retrieved",
		slog.Int64("product_id", product.ID),
		slog.String("name", product.Name),
	)
	return product, nil
}

func (ps *ProductService) GetByCategoryID(ctx context.Context, catID int64) ([]domain.Product, error) {
	const op = "service.product.GetByCategoryID"
	log := ps.log.With("op", op)

	products, err := ps.repo.ListByCategory(ctx, catID)
	switch {
	case errors.Is(err, repository.ErrCategoryNotFound):
		return nil, ErrCategoryNotFound
	case len(products) == 0:
		return nil, ErrNoProductsFound
	case err != nil:
		log.Error("repo error", slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("products retrieved", slog.Int64("category_id", catID), slog.Int("count", len(products)))
	return products, nil
}
