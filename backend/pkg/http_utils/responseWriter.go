package http_utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message" example:"Invalid request"`
}

// WriteError writes a JSON error response using the standard encoding/json package.
func WriteError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(ErrorResponse{
		Message: msg,
	}); err != nil {
		http.Error(w, "failed to encode error response", http.StatusInternalServerError)
		return
	}
}

// WriteJSON writes any data as JSON using the standard encoding/json package.
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode error response", http.StatusInternalServerError)
		return
	}
}
