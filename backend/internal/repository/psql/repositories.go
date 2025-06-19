package repository

import "github.com/jmoiron/sqlx"

type Repositories struct {
	ProductRepository  *ProductRepository
	CategoryRepository *CategoryRepository
}

func New(db *sqlx.DB) *Repositories {
	return &Repositories{
		ProductRepository:  NewProductRepository(db),
		CategoryRepository: NewCategoryRepository(db),
	}
}
