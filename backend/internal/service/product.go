package service

import (
	"context"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
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
	const op = "service.product.create"
	log := ps.log.With("op", op)
	if product == nil {
		log.Warn("attempted to create nil product")
		return 0, domain.ErrValidation("product is nil")
	}
	if product.Name == "" {
		log.Warn("validation failed: name is empty")
		return 0, domain.ErrValidation("product name is required")
	}
	id, err := ps.repo.CreateProduct(ctx, product)
	if err != nil {
		log.Error("failed to create product", slog.Any("error", err))
		return 0, err
	}
	log.Info("product created", "product_id", id)

	return id, nil
}
func (ps *ProductService) GetDetailed(ctx context.Context, id int64) (*domain.Product, error) {
	const op = "service.product.get_detailed"
	log := ps.log.With("op", op)
	var product *domain.Product
	product, err := ps.repo.GetProductWithAttrs(ctx, id)
	if err != nil {
		log.Error("failed to get product", slog.Any("product_id", id), slog.Any("error", err))
		return nil, err
	}
	log.Info("product retrieved", "product_id", product.ID)
	return product, nil
}
