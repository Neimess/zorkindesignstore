package httputils

import (
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/transport/dto"
	"github.com/mailru/easyjson"
)

func WriteError(w http.ResponseWriter, statusCode int, msg string) (int, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return easyjson.MarshalToWriter(dto.ErrorResponse{
		Message: msg,
	}, w)
}

func WriteJSON(w http.ResponseWriter, statusCode int, data easyjson.Marshaler) (int, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return easyjson.MarshalToWriter(data, w)
}
