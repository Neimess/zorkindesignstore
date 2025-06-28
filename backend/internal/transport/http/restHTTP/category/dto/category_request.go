package dto

// swagger:model CategoryCreateRequest
type CategoryCreateRequest struct {
	Name string `json:"name" example:"Плитка" validate:"required,min=2"`
}

// swagger:model CategoryUpdateRequest
type CategoryUpdateRequest struct {
	Name string `json:"name" example:"Керамогранит" validate:"required,min=2"`
}
