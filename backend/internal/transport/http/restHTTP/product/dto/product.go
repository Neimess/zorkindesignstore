package dto

import (
	"fmt"
	"time"

	ve "github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/go-playground/validator/v10"

	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	prodDom "github.com/Neimess/zorkin-store-project/internal/domain/product"
	serviceDom "github.com/Neimess/zorkin-store-project/internal/domain/service"
)

var validate = validator.New()

// ProductRequest описывает базовый payload для создания или обновления продукта.
// swagger:model ProductCreateRequest
type ProductRequest struct {
	Name        string                    `json:"name" example:"Керамогранит" validate:"required,min=2,max=255"`
	Price       float64                   `json:"price" example:"3490" validate:"required,gt=0"`
	Description *string                   `json:"description,omitempty" example:"Прочный плиточный материал" validate:"omitempty,min=2,max=1000"`
	CategoryID  int64                     `json:"category_id" example:"1" validate:"required,gt=0"`
	ImageURL    *string                   `json:"image_url,omitempty" example:"https://example.com/image.png" validate:"omitempty,url"`
	Attributes  []ProductAttributeRequest `json:"attributes,omitempty" validate:"omitempty,dive"` // required:false
	Services    []ProductServiceRequest   `json:"services,omitempty" validate:"omitempty,dive"`   // required:false
}

// ProductAttributeRequest описывает новый атрибут для создания.
// swagger:model ProductAttributeRequest
type ProductAttributeRequest struct {
	Name  string  `json:"name" example:"Объём" validate:"required,min=2,max=255"`
	Unit  *string `json:"unit,omitempty" example:"л" validate:"omitempty,min=1,max=50"`
	Value string  `json:"value" example:"1.25" validate:"required"`
}

// ProductServiceRequest описывает привязку услуги.
// swagger:model ProductServiceRequest
type ProductServiceRequest struct {
	ServiceID int64 `json:"service_id" example:"1" validate:"required,gt=0"`
}

// ProductResponse описывает ответ API при получении продукта.
// swagger:model ProductResponse
type ProductResponse struct {
	ProductID   int64                           `json:"product_id" example:"10"`
	Name        string                          `json:"name" example:"Керамогранит"`
	Price       float64                         `json:"price" example:"3490"`
	Description *string                         `json:"description,omitempty"`
	CategoryID  int64                           `json:"category_id" example:"1"`
	ImageURL    *string                         `json:"image_url,omitempty"`
	CreatedAt   time.Time                       `json:"created_at" example:"2025-06-20T15:00:00Z"`
	Attributes  []ProductAttributeValueResponse `json:"attributes,omitempty"`
	Services    []ProductServiceResponse        `json:"services,omitempty"`
}

// ProductAttributeValueResponse отвечает за элемент атрибута в ответе.
// swagger:model ProductAttributeValueResponse
type ProductAttributeValueResponse struct {
	AttributeID int64   `json:"attribute_id" example:"2"`
	Name        string  `json:"name" example:"Объём"`
	Unit        *string `json:"unit,omitempty" example:"л"`
	Value       string  `json:"value" example:"1.25"`
}

// ProductServiceResponse отвечает за элемент услуги в ответе.
// swagger:model ProductServiceResponse
type ProductServiceResponse struct {
	ID          int64   `json:"id" example:"1"`
	Name        string  `json:"name" example:"Монтаж"`
	Description *string `json:"description,omitempty" example:"Установка изделия"`
	Price       float64 `json:"price" example:"1500.00"`
}

// Validate проверяет поля ProductRequest.
func (r ProductRequest) Validate() error {
	var errs []ve.FieldError
	if err := validate.Struct(r); err != nil {
		if inv, ok := err.(*validator.InvalidValidationError); ok {
			return inv
		}
		for _, e := range err.(validator.ValidationErrors) {
			var msg string
			switch e.Field() {
			case "Name":
				msg = "name is required and must be 2-255 chars"
			case "Price":
				msg = "price is required and must be >0"
			case "CategoryID":
				msg = "category_id is required and must be >0"
			case "Description":
				msg = "description must be 2-1000 chars if set"
			case "ImageURL":
				msg = "image_url must be valid URL if set"
			case "Attributes":
				msg = "invalid existing attributes"
			case "Services":
				msg = "invalid services"
			default:
				msg = "invalid field"
			}
			errs = append(errs, ve.FieldError{Field: e.Field(), Message: msg})
		}
	}
	// вложенная валидация
	for i, a := range r.Attributes {
		if err := a.Validate(); err != nil {
			for _, fe := range err.(ve.ValidationErrorResponse).Errors {
				fe.Field = fmt.Sprintf("attributes[%d].%s", i, fe.Field)
				errs = append(errs, fe)
			}
		}
	}
	for i, s := range r.Services {
		if err := s.Validate(); err != nil {
			errs = append(errs, ve.FieldError{Field: fmt.Sprintf("services[%d]", i), Message: err.Error()})
		}
	}
	if len(errs) > 0 {
		return ve.ValidationErrorResponse{Errors: errs}
	}
	return nil
}

// Validate проверяет один ProductAttributeRequest.
func (r ProductAttributeRequest) Validate() error {
	var errs []ve.FieldError
	if err := validate.Struct(r); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var msg string
			switch e.Field() {
			case "Name":
				msg = "name required 2-255 chars"
			case "Unit":
				msg = "unit must be 1-50 chars if set"
			case "Value":
				msg = "value is required"
			}
			errs = append(errs, ve.FieldError{Field: e.Field(), Message: msg})
		}
	}
	if len(errs) > 0 {
		return ve.ValidationErrorResponse{Errors: errs}
	}
	return nil
}

// Validate проверяет один ProductServiceRequest.
func (r ProductServiceRequest) Validate() error {
	var errs []ve.FieldError
	if err := validate.Struct(r); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			if e.Field() == "ServiceID" {
				errs = append(errs, ve.FieldError{Field: "service_id", Message: "service_id must be >0"})
			}
		}
	}
	if len(errs) > 0 {
		return ve.ValidationErrorResponse{Errors: errs}
	}
	return nil
}

// MapCreateToDomain конвертирует ProductRequest в доменную модель.
func (r ProductRequest) MapCreateToDomain() *prodDom.Product {
	p := &prodDom.Product{
		Name:        r.Name,
		Price:       r.Price,
		CategoryID:  r.CategoryID,
		Description: r.Description,
		ImageURL:    r.ImageURL,
	}
	for _, a := range r.Attributes {
		p.Attributes = append(p.Attributes, prodDom.ProductAttribute{
			Attribute: attr.Attribute{
				Name:       a.Name,
				Unit:       a.Unit,
				CategoryID: r.CategoryID,
			},
			Value: a.Value,
		})
	}
	for _, s := range r.Services {
		p.Services = append(p.Services, serviceDom.Service{ID: s.ServiceID})
	}
	return p
}

// MapUpdateToDomain конвертирует ProductRequest в домен для обновления.
func (r *ProductRequest) MapUpdateToDomain(id int64) *prodDom.Product {
	p := r.MapCreateToDomain()
	p.ID = id
	return p
}

// MapDomainToProductResponse строит ProductResponse из доменной модели.
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
		resp.Attributes = append(resp.Attributes, ProductAttributeValueResponse{AttributeID: pa.AttributeID, Name: pa.Attribute.Name, Unit: pa.Attribute.Unit, Value: pa.Value})
	}
	for _, s := range p.Services {
		resp.Services = append(resp.Services, ProductServiceResponse{ID: s.ID, Name: s.Name, Description: s.Description, Price: s.Price})
	}
	return resp
}
