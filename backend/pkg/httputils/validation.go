// pkg/httputils/validation.go
package httputils

import (
	"encoding/json"
	"net/http"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Errors []FieldError `json:"errors"`
}

func (v ValidationErrorResponse) Error() string {
	return "validation failed"
}

func WriteValidationError(w http.ResponseWriter, statusCode int, resp ValidationErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}
