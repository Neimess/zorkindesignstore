package dto

import "time"

//swagger:model ProductCreateRequest
//go:generate easyjson -all
type ProductCreateRequest struct {
	Name        string                         `json:"name" example:"Керамогранит" validate:"required,min=2"`
	Price       *float64                       `json:"price" example:"3490" validate:"required,gt=0"`
	Description *string                        `json:"description,omitempty" example:"Прочный плиточный материал"`
	CategoryID  *int64                         `json:"category_id" example:"1" validate:"required,gt=0"`
	ImageURL    *string                        `json:"image_url,omitempty" example:"https://example.com/image.png" validate:"url"`
	Attributes  []ProductAttributeValueRequest `json:"attributes,omitempty" validate:"dive"`
}

// ProductAttributeValueRequest defines one attribute value in request.
//
//swagger:model ProductAttributeValueRequest
type ProductAttributeValueRequest struct {
	AttributeID *int64 `json:"attribute_id" example:"2" validate:"required,gt=1"`
	Value       string `json:"value" example:"1.25" validate:"required"`
}

// ProductResponse describes the created or fetched product.
//
//swagger:model ProductResponse
type ProductResponse struct {
	ProductID   int64                           `json:"product_id" example:"10"`
	Name        string                          `json:"name" example:"Керамогранит"`
	Price       float64                         `json:"price" example:"3490"`
	Description *string                         `json:"description,omitempty"`
	CategoryID  int64                           `json:"category_id" example:"1"`
	ImageURL    *string                         `json:"image_url,omitempty"`
	CreatedAt   time.Time                       `json:"created_at" example:"2025-06-20T15:00:00Z"`
	Attributes  []ProductAttributeValueResponse `json:"attributes,omitempty"`
}

//swagger:model ProductAttributeValueResponse
type ProductAttributeValueResponse struct {
	AttributeID  int64   `json:"attribute_id" example:"2"`
	Name         string  `json:"name" example:"Объём"`
	Slug         string  `json:"slug" example:"volume"`
	Unit         *string `json:"unit,omitempty" example:"л"`
	IsFilterable bool    `json:"is_filterable"`
	Value        string  `json:"value" example:"1.25"`
}
