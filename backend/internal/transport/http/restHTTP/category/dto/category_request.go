package dto

import (
	ve "github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

// swagger:model CategoryRequest
type CategoryRequest struct {
	Name string `json:"name" example:"Плитка" validate:"required,min=2,max=255"`
}

func (r CategoryRequest) Validate() error {
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
					Message: "name is required and must be between 2 and 255 characters",
				})
			default:
				errs = append(errs, ve.FieldError{
					Field:   e.Field(),
					Message: "invalid field",
				})
			}
		}
	}

	if len(errs) > 0 {
		return ve.ValidationErrorResponse{Errors: errs}
	}

	return nil
}
