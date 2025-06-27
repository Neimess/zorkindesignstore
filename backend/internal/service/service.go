// internal/service/deps.go

package service

import (
	"fmt"

	"github.com/Neimess/zorkin-store-project/internal/service/attribute"
	"github.com/Neimess/zorkin-store-project/internal/service/auth"
	"github.com/Neimess/zorkin-store-project/internal/service/category"
	"github.com/Neimess/zorkin-store-project/internal/service/preset"
	"github.com/Neimess/zorkin-store-project/internal/service/product"
)

type Deps struct {
	ProductRepo   product.ProductRepository
	CategoryRepo  category.CategoryRepository
	PresetRepo    preset.PresetRepository
	AttributeRepo attribute.AttributeRepository
	JWTGenerator  auth.JWTGenerator
}

func NewDeps(
	productRepo product.ProductRepository,
	categoryRepo category.CategoryRepository,
	presetRepo preset.PresetRepository,
	attributeRepo attribute.AttributeRepository,
	jwtGenerator auth.JWTGenerator,
) Deps {
	return Deps{
		ProductRepo:   productRepo,
		CategoryRepo:  categoryRepo,
		PresetRepo:    presetRepo,
		AttributeRepo: attributeRepo,
		JWTGenerator:  jwtGenerator,
	}
}

type Service struct {
	ProductService   *product.Service
	CategoryService  *category.Service
	AuthService      *auth.Service
	PresetService    *preset.Service
	AttributeService *attribute.Service
}

func New(d Deps) (*Service, error) {
	// Delegate dependency validation to service-specific NewDeps
	prodDeps, err := product.NewDeps(d.ProductRepo)
	if err != nil {
		return nil, fmt.Errorf("product service init: %w", err)
	}
	prodSvc := product.New(prodDeps)

	catDeps, err := category.NewDeps(d.CategoryRepo)
	if err != nil {
		return nil, fmt.Errorf("category service init: %w", err)
	}
	catSvc := category.New(catDeps)

	authDeps, err := auth.NewDeps(d.JWTGenerator)
	if err != nil {
		return nil, fmt.Errorf("auth service init: %w", err)
	}
	authSvc := auth.New(authDeps)

	presetDeps, err := preset.NewDeps(d.PresetRepo)
	if err != nil {
		return nil, fmt.Errorf("preset service init: %w", err)
	}
	presetSvc := preset.New(presetDeps)

	attrDeps, err := attribute.NewDeps(d.AttributeRepo, d.CategoryRepo)
	if err != nil {
		return nil, fmt.Errorf("attribute service init: %w", err)
	}
	attrSvc := attribute.New(attrDeps)

	return &Service{
		ProductService:   prodSvc,
		CategoryService:  catSvc,
		AuthService:      authSvc,
		PresetService:    presetSvc,
		AttributeService: attrSvc,
	}, nil
}
