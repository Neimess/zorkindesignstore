package dto

import (
	"fmt"

	ve "github.com/Neimess/zorkin-store-project/pkg/http_utils"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

//swaggo:model PresetRequest
type PresetRequest struct {
	Name        string              `json:"name" validate:"required,min=2,max=255"`
	Description *string             `json:"description,omitempty"`
	TotalPrice  float64             `json:"total_price" validate:"gt=0"`
	ImageURL    *string             `json:"image_url,omitempty" validate:"omitempty,url"`
	Items       []PresetRequestItem `json:"items" validate:"required,dive"`
}

func (r PresetRequest) Validate() error {
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
			case "TotalPrice":
				errs = append(errs, ve.FieldError{
					Field:   "total_price",
					Message: "total_price must be greater than 0",
				})
			case "ImageURL":
				errs = append(errs, ve.FieldError{
					Field:   "image_url",
					Message: "image_url must be a valid URL",
				})
			case "Items":
				errs = append(errs, ve.FieldError{
					Field:   "items",
					Message: "items must contain at least one item",
				})
			default:
				errs = append(errs, ve.FieldError{
					Field:   e.Field(),
					Message: "invalid field",
				})
			}
		}
	}

	for idx, item := range r.Items {
		if err := item.Validate(); err != nil {
			if veResp, ok := err.(ve.ValidationErrorResponse); ok {
				for _, ferr := range veResp.Errors {
					ferr.Field = fmt.Sprintf("items[%d].%s", idx, ferr.Field)
					errs = append(errs, ferr)
				}
			} else {
				errs = append(errs, ve.FieldError{
					Field:   fmt.Sprintf("items[%d]", idx),
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

//swaggo:model PresetRequestItem
type PresetRequestItem struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
}

func (i PresetRequestItem) Validate() error {

	var errs []ve.FieldError

	if err := validate.Struct(i); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			switch e.Field() {
			case "ProductID":
				errs = append(errs, ve.FieldError{
					Field:   "product_id",
					Message: "product_id is required",
				})
			case "Quantity":
				errs = append(errs, ve.FieldError{
					Field:   "quantity",
					Message: "quantity must be greater than 0",
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
