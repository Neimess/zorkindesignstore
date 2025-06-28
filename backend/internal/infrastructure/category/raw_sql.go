package category

import "github.com/Neimess/zorkin-store-project/internal/domain/category"

type categoryDB struct {
	ID   int64  `db:"category_id"`
	Name string `db:"name"`
}

func (r categoryDB) toDomain() *category.Category {
	return &category.Category{
		ID:   r.ID,
		Name: r.Name,
	}
}
