package dto

// swagger:model CategoryCreateRequest
//
//go:generate easyjson -all
type CategoryCreateRequest struct {
	Name string `json:"name" example:"Плитка" validate:"required,min=2"`
}

// swagger:model CategoryUpdateRequest
type CategoryUpdateRequest struct {
	Name string `json:"name" example:"Керамогранит" validate:"required,min=2"`
}

// swagger:model CategoryResponse
type CategoryResponse struct {
	ID   int64  `json:"id" example:"3"`
	Name string `json:"name" example:"Керамогранит"`
}

// swagger:model CategoryResponse
//
//easyjson:json
type CategoryResponseSlice []CategoryResponse
