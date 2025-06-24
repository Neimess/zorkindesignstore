package repository

import (
	"log/slog"
	"github.com/Neimess/zorkin-store-project/internal/infrastructure/category"
	"github.com/jmoiron/sqlx"
)

type Deps struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}

type Repositories struct {
	ProductRepository  *ProductRepository
	CategoryRepository *category.CategoryRepository
	PresetRepository   *PresetRepository
}

func New(deps Deps) *Repositories {
	r := &Repositories{
		ProductRepository:  NewProductRepository(deps.DB, deps.Logger),
		CategoryRepository: category.NewCategoryRepository(deps.DB),
		PresetRepository:   NewPresetRepository(deps.DB, deps.Logger),
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
	}
}
