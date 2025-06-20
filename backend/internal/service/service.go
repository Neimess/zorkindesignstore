// internal/service/deps.go

package service

import (
	"io"
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
		ProductService:  NewProductService(d.ProductRepository, silentLogger()),
		CategoryService: NewCategoryService(d.CategoryRepository, silentLogger()),
	}
}

func silentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
