package httputils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Validatable interface {
	Validate() error
}

func DecodeAndValidate[T Validatable](w http.ResponseWriter, r *http.Request, log *slog.Logger) (*T, bool) {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid JSON", slog.Any("error", err))
		WriteError(w, http.StatusBadRequest, "invalid JSON")
		return nil, false
	}
	if err := req.Validate(); err != nil {
		log.Warn("validation failed", slog.Any("error", err))
		if ve, ok := err.(ValidationErrorResponse); ok {
			WriteValidationError(w, http.StatusUnprocessableEntity, ve)
			return nil, false
		}
		WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return nil, false
	}
	return &req, true
}
