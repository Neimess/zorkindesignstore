package service

import (
	"context"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, p *domain.Product, images []domain.ProductImage) error
	GetProduct(ctx context.Context, id int64) (*domain.ProductWithImages, error)
}

type ProductService struct {
	repo ProductRepository
	log  *slog.Logger
}

func NewProductService(repo ProductRepository, logger *slog.Logger) *ProductService {
	return &ProductService{
		repo: repo,
		log:  logger.With("component", "product_service"),
	}
}

func (ps *ProductService) Create(ctx context.Context, product *domain.Product, images []domain.ProductImage) error {
	// TODO
	return nil
}
func (ps *ProductService) GetByID(ctx context.Context, id int64) (*domain.ProductWithImages, error) {
	// TODO
	return nil, nil
}
func (ps *ProductService) GetDetailed(ctx context.Context, id int64) (*domain.ProductWithCharacteristics, error) {
	// TODO
	return nil, nil
}
func (ps *ProductService) GetAll(ctx context.Context) ([]domain.ProductWithImages, error) {
	// TODO
	return nil, nil
}
