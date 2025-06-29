package dto

import (
	"time"

	prodDom "github.com/Neimess/zorkin-store-project/internal/domain/product"
)

//swagger:model ProductCreateRequest
type ProductCreateRequest struct {
	Name        string                         `json:"name" example:"Керамогранит" validate:"required,min=2"`
	Price       float64                        `json:"price" example:"3490" validate:"required,gt=0"`
	Description *string                        `json:"description,omitempty" example:"Прочный плиточный материал"`
	CategoryID  int64                          `json:"category_id" example:"1" validate:"required,gt=0"`
	ImageURL    *string                        `json:"image_url,omitempty" example:"https://example.com/image.png" validate:"url"`
	Attributes  []ProductAttributeValueRequest `json:"attributes,omitempty" validate:"dive"`
}

type ProductUpdateRequest struct {
	Name        string                         `json:"name"        validate:"required"`
	Price       float64                        `json:"price"       validate:"required,gt=0"`
	Description *string                        `json:"description,omitempty" validate:"omitempty"`
	CategoryID  int64                          `json:"category_id" validate:"required,gt=0"`
	ImageURL    *string                        `json:"image_url,omitempty" validate:"omitempty,url"`
	Attributes  []ProductAttributeValueRequest `json:"attributes,omitempty" validate:"omitempty,dive"`
}

//swagger:model ProductAttributeValueRequest
type ProductAttributeValueRequest struct {
	AttributeID int64  `json:"attribute_id" example:"2" validate:"required,gt=1"`
	Value       string `json:"value" example:"1.25" validate:"required"`
}

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
	AttributeID int64   `json:"attribute_id" example:"2"`
	Name        string  `json:"name" example:"Объём"`
	Unit        *string `json:"unit,omitempty" example:"л"`
	Value       string  `json:"value" example:"1.25"`
}

func MapCreateReqToDomain(req *ProductCreateRequest) *prodDom.Product {
	p := &prodDom.Product{
		Name:        req.Name,
		Price:       req.Price,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		ImageURL:    req.ImageURL,
	}

	for _, a := range req.Attributes {
		p.Attributes = append(p.Attributes, prodDom.ProductAttribute{
			AttributeID: a.AttributeID,
			Value:       a.Value,
		})
	}
	return p
}

func MapDomainToProductResponse(p *prodDom.Product) *ProductResponse {
	resp := &ProductResponse{
		ProductID:   p.ID,
		Name:        p.Name,
		Price:       p.Price,
		CategoryID:  p.CategoryID,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt,
	}
	for _, pa := range p.Attributes {

		resp.Attributes = append(resp.Attributes, ProductAttributeValueResponse{
			AttributeID: pa.AttributeID,
			Name:        pa.Attribute.Name,
			Unit:        pa.Attribute.Unit,
			Value:       pa.Value,
		})
	}
	return resp
}

func MapUpdateReqToDomain(id int64, req *ProductUpdateRequest) *prodDom.Product {
	p := &prodDom.Product{
		ID:          id,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		CategoryID:  req.CategoryID,
	}

	if req.Attributes != nil {
		for _, a := range req.Attributes {
			p.Attributes = append(p.Attributes, prodDom.ProductAttribute{
				AttributeID: a.AttributeID,
				Value:       a.Value,
			})
		}
	}
	return p
}
