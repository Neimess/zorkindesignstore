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
	Logger                      *slog.Logger
	Config                      *config.Config
	JWTGenerator                JWTGenerator
	ProductRepository           ProductRepository
	CategoryRepository          CategoryRepository
	CategoryAttributeRepository CategoryAttributeRepository
	PresetRepository            PresetRepository
}

type Service struct {
	ProductService           *ProductService
	CategoryService          *CategoryService
	AuthService              *AuthService
	CategoryAttributeService *CategoryAttributeService
	PresetService            *PresetService
}

func New(d Deps) *Service {
	return &Service{
		ProductService:           NewProductService(d.ProductRepository, d.Logger),
		CategoryService:          NewCategoryService(d.CategoryRepository, d.Logger),
		AuthService:              NewAuthService(d.JWTGenerator, d.Logger),
		CategoryAttributeService: NewCategoryAttributeService(d.CategoryAttributeRepository, d.Logger),
		PresetService:            NewPresetService(d.PresetRepository, d.Logger),
	}
}

func silentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
