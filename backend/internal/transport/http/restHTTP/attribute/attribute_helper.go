package attribute

import (
	"errors"
	"log/slog"
	"net/http"

	attrDom "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
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
	// атрибут не найден
	case errors.Is(err, attrDom.ErrAttributeNotFound):
		httputils.WriteError(w, http.StatusNotFound, "attribute not found")

	// конфликт уникальности имени атрибута
	case errors.Is(err, attrDom.ErrAttributeAlreadyExists):
		httputils.WriteError(w, http.StatusConflict, "attribute already exists")

	// категория не найдена
	case errors.Is(err, catDom.ErrCategoryNotFound):
		httputils.WriteError(w, http.StatusNotFound, "category not found")

	// пустой батч
	case errors.Is(err, attrDom.ErrBatchEmpty):
		httputils.WriteError(w, http.StatusBadRequest, "no attributes provided for batch")

	// всё прочее — 500
	default:
		h.log.Error("service error", slog.Any("error", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
