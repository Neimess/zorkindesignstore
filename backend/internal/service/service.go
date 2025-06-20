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

func New(d Deps) (*Service, error) {
	ps, err := NewProductService(d.ProductRepository, d.Logger)
	if err != nil {
		return nil, err
	}
	// cs, err := NewCategoryService(d.CategoryRepository, d.Logger)
	// if err != nil {
	// 	return nil, err
	// }
	return &Service{
		ProductService: ps,
		// CategoryService: cs,
	}, nil
}

func silentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
