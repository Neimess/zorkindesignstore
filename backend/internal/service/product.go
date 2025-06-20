package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	repository "github.com/Neimess/zorkin-store-project/internal/repository/psql"
)

var (
	ErrProductNameEmpty = errors.New("product name is empty")
	ErrProductIsNil     = errors.New("product is nil")
	ErrProductNotFound  = errors.New("product not found")
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, p *domain.Product) (int64, error)
	GetProductWithAttrs(ctx context.Context, id int64) (*domain.Product, error)
}

type ProductService struct {
	repo ProductRepository
	log  *slog.Logger
}

func NewProductService(repo ProductRepository, logger *slog.Logger) *ProductService {
	return &ProductService{
		repo: repo,
		log:  logger,
	}
}

func (ps *ProductService) Create(ctx context.Context, product *domain.Product) (int64, error) {
	const op = "service.product.Create"
	log := ps.log.With("op", op)
	if product == nil {
		log.Warn("attempted to create nil product")
		return 0, ErrProductIsNil
	}
	if product.Name == "" {
		log.Warn("validation failed: name is empty")
		return 0, ErrProductNameEmpty
	}
	id, err := ps.repo.CreateProduct(ctx, product)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn("referenced category does not exist",
				slog.Int("category_id", int(product.CategoryID)),
			)
		}
		log.Error("failed to create product", slog.Any("error", err))
		return 0, err
	}
	ps.log.Info("product created",
		slog.Int64("product_id", id),
		slog.String("name", product.Name),
		slog.Int("category_id", int(product.CategoryID)),
	)

	return id, nil
}
func (ps *ProductService) GetDetailed(ctx context.Context, id int64) (*domain.Product, error) {
	const op = "service.product.GetDetailed"
	log := ps.log.With("op", op)
	product, err := ps.repo.GetProductWithAttrs(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Info("not found product", slog.Any("product_id", id), slog.Any("error", err))
		}
		log.Info("failed to get product", slog.Any("product_id", id), slog.Any("error", err))
		return nil, err
	}
	log.Info("product retrieved",
		slog.Int64("product_id", product.ID),
		slog.String("name", product.Name),
	)
	return product, nil
}
