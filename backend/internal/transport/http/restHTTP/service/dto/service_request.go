package dto

import (
	ve "github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

type ServiceRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description,omitempty"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

func (r ServiceRequest) Validate() error {
	var errs []ve.FieldError
	if err := validate.Struct(r); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			switch e.Field() {
			case "Name":
				errs = append(errs, ve.FieldError{Field: "name", Message: "name is required and must be 1-255 chars"})
			case "Price":
				errs = append(errs, ve.FieldError{Field: "price", Message: "price is required and must be > 0"})
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
