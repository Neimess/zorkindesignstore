package attribute

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
)

func (h *Handler) parseCategoryID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := httputils.IDFromURL(r, "categoryID")
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid category id")
		return 0, false
	}
	return id, true
}

func (h *Handler) parseAttributeID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid attribute id")
		return 0, false
	}
	return id, true
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrAttributeNameEmpty):
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		httputils.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrConflict):
		httputils.WriteError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrCategoryNotFound):
		httputils.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInternal):
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
	case errors.Is(err, domain.ErrBadRequest):
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrValidation):
		httputils.WriteError(w, http.StatusUnprocessableEntity, "validation error")
	default:
		h.log.Error("service error", slog.Any("error", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal error")
	}
}
