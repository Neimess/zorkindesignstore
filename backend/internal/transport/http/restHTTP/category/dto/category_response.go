package dto

// swagger:model CategoryResponse
type CategoryResponse struct {
	ID       int64  `json:"id" example:"3"`
	Name     string `json:"name" example:"Керамогранит"`
	ParentID *int64 `json:"parent_id,omitempty" example:"1"`
}
