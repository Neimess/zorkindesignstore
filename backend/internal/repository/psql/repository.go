package repository

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Deps struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}

type Repositories struct {
	ProductRepository           *ProductRepository
	CategoryRepository          *CategoryRepository
	CategoryAttributeRepository *CategoryAttributeRepository
	PresetRepository            *PresetRepository
}

func New(deps Deps) *Repositories {
	r := &Repositories{
		ProductRepository:           NewProductRepository(deps.DB, deps.Logger),
		CategoryRepository:          NewCategoryRepository(deps.DB, deps.Logger),
		CategoryAttributeRepository: NewCategoryAttributeRepository(deps.DB, deps.Logger),
		PresetRepository:            NewPresetRepository(deps.DB, deps.Logger),
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
	case r.CategoryAttributeRepository == nil:
		panic("CategoryAttributeRepository is not initialized")
	case r.PresetRepository == nil:
		panic("PresetRepository is not initialized")
	}
}
