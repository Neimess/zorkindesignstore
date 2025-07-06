package dto

import (
	"github.com/Neimess/zorkin-store-project/internal/domain/category"
)

func (in *CategoryRequest) ToDomainCreate() *category.Category {
	return &category.Category{
		Name:     in.Name,
		ParentID: in.ParentID,
	}
}

func (in *CategoryRequest) ToDomainUpdate(category_id int64) *category.Category {
	return &category.Category{
		ID:       category_id,
		Name:     in.Name,
		ParentID: in.ParentID,
	}
}

func ToDTOResponse(cat *category.Category) CategoryResponse {
	return CategoryResponse{
		ID:       cat.ID,
		Name:     cat.Name,
		ParentID: cat.ParentID,
	}
}

func ToDTOList(cats []category.Category, total int64) []CategoryResponse {
	out := make([]CategoryResponse, len(cats))
	for i, c := range cats {
		out[i] = ToDTOResponse(&c)
	}
	return out
}
