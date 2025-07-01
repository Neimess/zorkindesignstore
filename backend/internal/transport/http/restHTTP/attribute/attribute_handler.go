package attribute

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"log/slog"

	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/attribute/dto"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
)

type AttributeService interface {
	CreateAttributesBatch(ctx context.Context, category_id int64, input []attr.Attribute) error
	CreateAttribute(ctx context.Context, category_id int64, in *attr.Attribute) (*attr.Attribute, error)
	GetAttribute(ctx context.Context, categoryID, id int64) (*attr.Attribute, error)
	ListAttributes(ctx context.Context, categoryID int64) ([]attr.Attribute, error)
	UpdateAttribute(ctx context.Context, attr *attr.Attribute) (*attr.Attribute, error)
	DeleteAttribute(ctx context.Context, id int64) error
}

type Deps struct {
	srv AttributeService
	log *slog.Logger
}

func NewDeps(log *slog.Logger, srv AttributeService) (Deps, error) {
	if srv == nil {
		return Deps{}, errors.New("attribute handler: missing AttributeService")
	}
	if log == nil {
		return Deps{}, errors.New("attribute handler: missing Logger")
	}
	return Deps{
		srv: srv,
		log: log.With("component", "rest.attribute"),
	}, nil
}

type Handler struct {
	srv AttributeService
	log *slog.Logger
}

func New(d Deps) *Handler {

	return &Handler{
		srv: d.srv,
		log: d.log,
	}
}

// CreateAttributesBatch godoc
// @Summary      Batch create attributes
// @Description  Создать несколько атрибутов в одной категории
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        categoryID path int true "Category ID"
// @Param        data body dto.CreateAttributesBatchRequest true "Batch input"
// @Success      204 "No Content"
// @Failure      400 {object} httputils.ErrorResponse
// @Failure      422 {object} httputils.ErrorResponse
// @Failure      500 {object} httputils.ErrorResponse
// @Router       /api/admin/category/{categoryID}/attribute/batch [post]
func (h *Handler) CreateAttributesBatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID, ok := h.parseCategoryID(w, r)
	if !ok {
		return
	}

	var req dto.CreateAttributesBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("invalid JSON", slog.Any("error", err))
		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := req.Validate(); err != nil {
		h.log.Warn("validation failed", slog.Any("error", err))
		if ve, ok := err.(httputils.ValidationErrorResponse); ok {
			httputils.WriteValidationError(w, http.StatusUnprocessableEntity, ve)
			return
		}
		httputils.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	attrs := req.MapToDomainBatch()
	if err := h.srv.CreateAttributesBatch(ctx, categoryID, attrs); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateAttribute godoc
// @Summary      Create an attribute
// @Description  Создать один атрибут в категории
// @Security     BearerAuth
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        categoryID path int true "Category ID"
// @Param        data body dto.AttributeRequest true "Attribute data"
// @Success      201 "Created"
// @Failure      400 {object} httputils.ErrorResponse
// @Failure      422 {object} httputils.ErrorResponse
// @Failure      500 {object} httputils.ErrorResponse
// @Router       /api/admin/category/{categoryID}/attribute [post]
func (h *Handler) CreateAttribute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID, ok := h.parseCategoryID(w, r)
	if !ok {
		return
	}

	var req dto.AttributeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("invalid JSON", slog.Any("error", err))
		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := req.Validate(); err != nil {
		h.log.Warn("validation failed", slog.Any("error", err))
		if ve, ok := err.(httputils.ValidationErrorResponse); ok {
			httputils.WriteValidationError(w, http.StatusUnprocessableEntity, ve)
			return
		}
		httputils.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	attr := req.MapToDomain()
	attr.CategoryID = categoryID
	created, err := h.srv.CreateAttribute(ctx, categoryID, attr)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/api/admin/category/%d/attribute/%d", categoryID, created.ID))
	httputils.WriteJSON(w, http.StatusCreated, dto.MapToAttributeResponse(created))
}

// ListAttributes godoc
// @Summary      Get all attributes for a category
// @Description  Получить все атрибуты в категории
// @Tags         categories
// @Param        categoryID path int true "Category ID"
// @Produce      json
// @Success      200 {array}  dto.AttributeResponse
// @Failure      400 {object} httputils.ErrorResponse
// @Failure      404 {object} httputils.ErrorResponse
// @Failure      500 {object} httputils.ErrorResponse
// @Router       /api/category/{categoryID}/attribute [get]
func (h *Handler) ListAttributes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID, ok := h.parseCategoryID(w, r)
	if !ok {
		return
	}

	attrs, err := h.srv.ListAttributes(ctx, categoryID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	resp := make(dto.AttributeListResponse, len(attrs))
	for i, a := range attrs {
		resp[i] = dto.AttributeResponse{
			ID:         a.ID,
			Name:       a.Name,
			Unit:       a.Unit,
			CategoryID: a.CategoryID,
		}
	}
	httputils.WriteJSON(w, http.StatusOK, resp)
}

// GetAttribute godoc
// @Summary      Get single attribute
// @Description  Получить один атрибут по ID и категории
// @Tags         categories
// @Param        categoryID path int true "Category ID"
// @Param        id path int true "Attribute ID"
// @Produce      json
// @Success      200 {object} dto.AttributeResponse
// @Failure      400 {object} httputils.ErrorResponse
// @Failure      404 {object} httputils.ErrorResponse
// @Failure      500 {object} httputils.ErrorResponse
// @Router       /api/category/{categoryID}/attribute/{id} [get]
func (h *Handler) GetAttribute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID, ok := h.parseCategoryID(w, r)
	if !ok {
		return
	}
	id, ok := h.parseAttributeID(w, r)
	if !ok {
		return
	}

	attr, err := h.srv.GetAttribute(ctx, categoryID, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	resp := dto.AttributeResponse{
		ID:         attr.ID,
		Name:       attr.Name,
		Unit:       attr.Unit,
		CategoryID: attr.CategoryID,
	}
	httputils.WriteJSON(w, http.StatusOK, resp)
}

// UpdateAttribute godoc
// @Summary      Update an attribute
// @Description  Обновить атрибут по ID и категории
// @Tags         categories
// @Security     BearerAuth
// @Accept       json
// @Param        categoryID path int true "Category ID"
// @Param        id path int true "Attribute ID"
// @Param        data body dto.AttributeRequest true "Attribute data"
// @Success      200 {object} dto.AttributeResponse
// @Failure      400 {object} httputils.ErrorResponse
// @Failure      422 {object} httputils.ErrorResponse
// @Failure      500 {object} httputils.ErrorResponse
// @Router       /api/admin/category/{categoryID}/attribute/{id} [put]
func (h *Handler) UpdateAttribute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID, ok := h.parseCategoryID(w, r)
	if !ok {
		return
	}
	id, ok := h.parseAttributeID(w, r)
	if !ok {
		return
	}

	var req dto.AttributeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := req.Validate(); err != nil {
		h.log.Warn("validation failed", slog.Any("error", err))
		if ve, ok := err.(httputils.ValidationErrorResponse); ok {
			httputils.WriteValidationError(w, http.StatusUnprocessableEntity, ve)
			return
		}
		httputils.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	attr := req.MapToDomain()

	attr.ID = id
	attr.CategoryID = categoryID

	updated, err := h.srv.UpdateAttribute(ctx, attr)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	httputils.WriteJSON(w, http.StatusOK, dto.MapToAttributeResponse(updated))
}

// DeleteAttribute godoc
// @Summary      Delete an attribute
// @Description  Удалить атрибут по ID
// @Tags         categories
// @Security     BearerAuth
// @Param        categoryID path int true "Category ID"
// @Param        id path int true "Attribute ID"
// @Success      204 "No Content"
// @Failure      400 {object} httputils.ErrorResponse
// @Failure      404 {object} httputils.ErrorResponse
// @Failure      500 {object} httputils.ErrorResponse
// @Router       /api/admin/category/{categoryID}/attribute/{id} [delete]
func (h *Handler) DeleteAttribute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := h.parseAttributeID(w, r)
	if !ok {
		return
	}

	if err := h.srv.DeleteAttribute(ctx, id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
