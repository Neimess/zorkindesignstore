package repository

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	ProductRepository  *ProductRepository
	CategoryRepository *CategoryRepository
}

func New(db *sqlx.DB, logger *slog.Logger) *Repositories {
	return &Repositories{
		ProductRepository:  NewProductRepository(db, logger),
		// CategoryRepository: NewCategoryRepository(db, logger),
	}
}
