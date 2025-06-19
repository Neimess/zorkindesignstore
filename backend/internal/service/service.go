// internal/service/deps.go

package service

import (
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/config"
)

type Deps struct {
	Logger             *slog.Logger
	Config             *config.Config
	ProductRepository  ProductRepository
	CategoryRepository CategoryRepository
}

type Service struct {
	ProductService  *ProductService
	CategoryService *CategoryService
}

func New(d Deps) *Service {
	return &Service{
		ProductService:  NewProductService(d.ProductRepository, d.Logger),
		CategoryService: NewCategoryService(d.CategoryRepository, d.Logger),
	}
}

