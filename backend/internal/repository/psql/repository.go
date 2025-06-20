package repository

import (
	"errors"
)

var (
	ErrNotFound = errors.New("target not found")
)

type Repositories struct {
	ProductRepository  *ProductRepository
	CategoryRepository *CategoryRepository
}

func New(deps Deps) *Repositories {
	return &Repositories{
		ProductRepository: NewProductRepository(deps.DB, deps.Logger),
		// CategoryRepository: NewCategoryRepository(db, logger),
	}
}
