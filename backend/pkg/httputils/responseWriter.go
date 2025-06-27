package httputils

import (
	"encoding/json"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/transport/dto"
)

// WriteError writes a JSON error response using the standard encoding/json package.
func WriteError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// dto.ErrorResponse is something like:
	//   type ErrorResponse struct { Message string `json:"message"` }
	if err := json.NewEncoder(w).Encode(dto.ErrorResponse{
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
