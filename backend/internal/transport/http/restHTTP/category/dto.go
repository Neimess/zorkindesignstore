package category

// swagger:model CategoryCreateRequest
//
//go:generate easyjson -all
type CategoryCreateRequest struct {
	Name string `json:"name" example:"Плитка" validate:"required,min=2"`
}

// swagger:model CategoryUpdateRequest
//
//easyjson:json
type CategoryUpdateRequest struct {
	Name string `json:"name" example:"Керамогранит" validate:"required,min=2"`
}

// swagger:model CategoryResponse
//
//easyjson:json
type CategoryResponse struct {
	ID   int64  `json:"id" example:"3"`
	Name string `json:"name" example:"Керамогранит"`
	Attributes []AttributePayload `json:"attributes,omitempty"`
}


//swagger:model CreateCategoryReq
//
//easyjson:json
type CreateCategoryReq struct {
	Name       string `json:"name"`
	Attributes []AttributePayload `json:"attributes"`
}

// swagger:model AttributePayload
//
//easyjson:json
type AttributePayload struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Unit         string `json:"unit,omitempty"`
	IsFilterable bool   `json:"is_filterable"`
	IsRequired   bool   `json:"is_required"`
	Priority     int    `json:"priority"`
}

// swagger:model CategoryResponseList
//
//easyjson:json
type CategoryResponseList struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int64              `json:"total"`
}