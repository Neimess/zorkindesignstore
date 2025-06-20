// internal/service/deps.go

package service

import (
	"io"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/config"
)

type JWTGenerator interface {
	Generate(userID string) (string, error)
}

type Deps struct {
	Logger             *slog.Logger
	Config             *config.Config
	JWTGenerator       JWTGenerator
	ProductRepository  ProductRepository
	CategoryRepository CategoryRepository
}

type Service struct {
	ProductService  *ProductService
	CategoryService *CategoryService
	AuthService     *AuthService
}

func New(d Deps) *Service {
	return &Service{
		ProductService:  NewProductService(d.ProductRepository, d.Logger),
		CategoryService: NewCategoryService(d.CategoryRepository, d.Logger),
		AuthService:     NewAuthService(d.JWTGenerator, d.Logger),
	}
}

func silentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
