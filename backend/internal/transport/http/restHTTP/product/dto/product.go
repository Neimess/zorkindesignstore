package dto

import (
	"fmt"
	"time"

	ve "github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/go-playground/validator/v10"

	prodDom "github.com/Neimess/zorkin-store-project/internal/domain/product"
)

var validate *validator.Validate = validator.New()

//swagger:model ProductCreateRequest
type ProductRequest struct {
	Name        string                         `json:"name" example:"Керамогранит" validate:"required,min=2"`
	Price       float64                        `json:"price" example:"3490" validate:"required,gt=0"`
	Description *string                        `json:"description,omitempty" example:"Прочный плиточный материал"`
	CategoryID  int64                          `json:"category_id" example:"1" validate:"required,gt=0"`
	ImageURL    *string                        `json:"image_url,omitempty" example:"https://example.com/image.png" validate:"omitempty,url"`
	Attributes  []ProductAttributeValueRequest `json:"attributes,omitempty" validate:"dive"`
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

func (r ProductRequest) Validate() error {
	var errs []ve.FieldError
	if err := validate.Struct(r); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			switch e.Field() {
			case "Name":
				errs = append(errs, ve.FieldError{Field: "name", Message: "name is required and must be at least 2 characters"})
			case "Price":
				errs = append(errs, ve.FieldError{Field: "price", Message: "price is required and must be greater than 0"})
			case "CategoryID":
				errs = append(errs, ve.FieldError{Field: "category_id", Message: "category_id is required and must be greater than 0"})
			case "ImageURL":
				errs = append(errs, ve.FieldError{Field: "image_url", Message: "image_url must be a valid URL"})
			case "Attributes":
				errs = append(errs, ve.FieldError{Field: "attributes", Message: "invalid attributes"})
			default:
				errs = append(errs, ve.FieldError{Field: e.Field(), Message: "invalid field"})
			}
		}
	}
	for idx, attr := range r.Attributes {
		if err := attr.Validate(); err != nil {
			if veResp, ok := err.(ve.ValidationErrorResponse); ok {
				for _, ferr := range veResp.Errors {
					ferr.Field = fmt.Sprintf("attributes[%d].%s", idx, ferr.Field)
					errs = append(errs, ferr)
				}
			} else {
				errs = append(errs, ve.FieldError{Field: fmt.Sprintf("attributes[%d]", idx), Message: err.Error()})
			}
		}
	}
	if len(errs) > 0 {
		return ve.ValidationErrorResponse{Errors: errs}
	}
	return nil
}

func (r ProductAttributeValueRequest) Validate() error {
	var errs []ve.FieldError
	if err := validate.Struct(r); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			switch e.Field() {
			case "AttributeID":
				errs = append(errs, ve.FieldError{Field: "attribute_id", Message: "attribute_id is required and must be greater than 1"})
			case "Value":
				errs = append(errs, ve.FieldError{Field: "value", Message: "value is required"})
			default:
				errs = append(errs, ve.FieldError{Field: e.Field(), Message: "invalid field"})
			}
		}
	}
	if len(errs) > 0 {
		return ve.ValidationErrorResponse{Errors: errs}
	}
	return nil
}

func MapCreateReqToDomain(req *ProductRequest) *prodDom.Product {
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

func MapUpdateReqToDomain(id int64, req *ProductRequest) *prodDom.Product {
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
