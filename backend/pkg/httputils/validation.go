// pkg/httputils/validation.go
package httputils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// FieldError describes a validation error on a single field
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ValidationErrorResponse is the payload returned on validation failure
type ValidationErrorResponse struct {
	Errors []FieldError `json:"errors"`
}

// RespondValidationError writes a 422 Unprocessable Entity response
// with details about each failed field validation.
func RespondValidationError(w http.ResponseWriter, ve validator.ValidationErrors) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)

	out := make([]FieldError, 0, len(ve))
	for _, fe := range ve {
		out = append(out, FieldError{
			Field: fe.Field(),
			Error: fe.Error(), // or fmt.Sprintf("failed on '%s' tag", fe.Tag())
		})
	}

	if err := json.NewEncoder(w).Encode(ValidationErrorResponse{Errors: out}); err != nil {
		http.Error(w, "failed to encode validation error response", http.StatusInternalServerError)
		return
	}
}
