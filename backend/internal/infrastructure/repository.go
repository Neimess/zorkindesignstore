package repository

import (
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/infrastructure/attribute"
	"github.com/Neimess/zorkin-store-project/internal/infrastructure/category"
	"github.com/Neimess/zorkin-store-project/internal/infrastructure/preset"
	"github.com/Neimess/zorkin-store-project/internal/infrastructure/product"
	"github.com/jmoiron/sqlx"
)

type Deps struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}

type Repositories struct {
	ProductRepository   *product.PGProductRepository
	CategoryRepository  *category.PGCategoryRepository
	PresetRepository    *preset.PGPresetRepository
	AttributeRepository *attribute.PGAttributeRepository
}

func New(deps Deps) *Repositories {
	r := &Repositories{
		ProductRepository:   product.NewPGProductRepository(deps.DB, deps.Logger),
		CategoryRepository:  category.NewPGCategoryRepository(deps.DB),
		PresetRepository:    preset.NewPGPresetRepository(deps.DB, deps.Logger),
		AttributeRepository: attribute.NewPGAttributeRepository(deps.DB),
	}

	r.mustValidate()
	return r
}

func (r *Repositories) mustValidate() {
	switch {
	case r.ProductRepository == nil:
		panic("ProductRepository is not initialized")
	case r.CategoryRepository == nil:
		panic("CategoryRepository is not initialized")
	case r.PresetRepository == nil:
		panic("PresetRepository is not initialized")
	case r.AttributeRepository == nil:
		panic("AttributeRepository is not initialized")
	}
}
