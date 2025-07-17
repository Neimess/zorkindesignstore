package dto

import (
	"fmt"

	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	ve "github.com/Neimess/zorkin-store-project/pkg/http_utils"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

type AttributeRequest struct {
	Name string  `json:"name" validate:"required,min=1,max=255"`
	Unit *string `json:"unit,omitempty" validate:"omitempty,max=50"`
}

func (r AttributeRequest) Validate() error {
	var errs []ve.FieldError
	if err := validate.Struct(r); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			switch e.Field() {
			case "Name":
				errs = append(errs, ve.FieldError{
					Field:   "name",
					Message: "name is required and must be between 1 and 255 characters",
				})
			case "Unit":
				errs = append(errs, ve.FieldError{
					Field:   "unit",
					Message: "unit must not exceed 50 characters",
				})
			default:
				errs = append(errs, ve.FieldError{
					Field:   e.Field(),
					Message: "invalid field",
				})
			}
		}
	}
	if r.Unit != nil && *r.Unit == "" {
		r.Unit = nil
	}
	if len(errs) > 0 {
		return ve.ValidationErrorResponse{Errors: errs}
	}
	return nil
}

type CreateAttributesBatchRequest struct {
	Items []AttributeRequest `json:"data" validate:"required,min=1,dive"`
}

func (r CreateAttributesBatchRequest) Validate() error {
	var errs []ve.FieldError
	if err := validate.Struct(r); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			errs = append(errs, ve.FieldError{
				Field:   e.Field(),
				Message: e.Error(),
			})
		}
	}
	for idx, item := range r.Items {
		if err := item.Validate(); err != nil {
			if veResp, ok := err.(ve.ValidationErrorResponse); ok {
				for _, ferr := range veResp.Errors {
					ferr.Field = fmt.Sprintf("data[%d].%s", idx, ferr.Field)
					errs = append(errs, ferr)
				}
			} else {
				errs = append(errs, ve.FieldError{
					Field:   fmt.Sprintf("data[%d]", idx),
					Message: err.Error(),
				})
			}
		}
	}
	if len(errs) > 0 {
		return ve.ValidationErrorResponse{Errors: errs}
	}
	return nil
}

func (r AttributeRequest) MapToDomain() *attr.Attribute {
	return &attr.Attribute{
		Name: r.Name,
		Unit: r.Unit,
	}
}

func (r CreateAttributesBatchRequest) MapToDomainBatch() []attr.Attribute {
	out := make([]attr.Attribute, len(r.Items))
	for i, it := range r.Items {
		out[i] = attr.Attribute{
			Name: it.Name,
			Unit: it.Unit,
		}
	}
	return out
}
