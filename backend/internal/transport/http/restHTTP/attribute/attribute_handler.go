package attribute

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"log/slog"

	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	"github.com/Neimess/zorkin-store-project/internal/service/attribute"
	_ "github.com/Neimess/zorkin-store-project/internal/transport/dto"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
	cv "github.com/Neimess/zorkin-store-project/pkg/validator"
	"github.com/go-playground/validator/v10"
)

type AttributeService interface {
	CreateAttributesBatch(ctx context.Context, input attribute.CreateAttributesBatchInput) error
	CreateAttribute(ctx context.Context, in *attribute.CreateAttributeInput) (*attr.Attribute, error)
	GetAttribute(ctx context.Context, categoryID, id int64) (*attr.Attribute, error)
	ListAttributes(ctx context.Context, categoryID int64) ([]attr.Attribute, error)
	UpdateAttribute(ctx context.Context, attr *attr.Attribute) error
	DeleteAttribute(ctx context.Context, id int64) error
}

type Deps struct {
	srv AttributeService
}

func NewDeps(srv AttributeService) (*Deps, error) {
	if srv == nil {
		return nil, errors.New("attribute handler: missing AttributeService")
	}
	return &Deps{
		srv: srv,
	}, nil
}

type Handler struct {
	srv      AttributeService
	log      *slog.Logger
	validate *validator.Validate
}

func New(d *Deps) *Handler {

	return &Handler{
		srv:      d.srv,
		log:      slog.Default().With("component", "rest.attribute"),
		validate: cv.GetValidator(),
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
// @Param        data body attribute.CreateAttributesBatchInput true "Batch input"
// @Success      204 "No Content"
// @Failure      400 {object} dto.ErrorResponse
// @Failure      422 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/admin/category/{categoryID}/attribute/batch [post]
func (h *Handler) CreateAttributesBatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID, ok := h.parseCategoryID(w, r)
	if !ok {
		return
	}

	var req attribute.CreateAttributesBatchInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("invalid JSON", slog.Any("error", err))
		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	req[0].CategoryID = categoryID
	if err := h.validate.StructCtx(ctx, &req); err != nil {
		h.log.Warn("validation failed", slog.Any("error", err))
		if ve, ok := err.(validator.ValidationErrors); ok {
			httputils.RespondValidationError(w, ve)
		} else {
			httputils.WriteError(w, http.StatusUnprocessableEntity, "invalid attribute data")
		}
		return
	}

	if err := h.srv.CreateAttributesBatch(ctx, req); err != nil {
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
// @Param        data body CreateAttributeReq true "Attribute data"
// @Success      201 {object} AttributeResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      422 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/admin/category/{categoryID}/attribute [post]
func (h *Handler) CreateAttribute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categoryID, ok := h.parseCategoryID(w, r)
	if !ok {
		return
	}

	var req CreateAttributeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("invalid JSON", slog.Any("error", err))
		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := h.validate.StructCtx(ctx, &req); err != nil {
		h.log.Warn("validation failed", slog.Any("error", err))
		if ve, ok := err.(validator.ValidationErrors); ok {
			httputils.RespondValidationError(w, ve)
		} else {
			httputils.WriteError(w, http.StatusUnprocessableEntity, "invalid attribute data")
		}
		return
	}

	attr, err := h.srv.CreateAttribute(ctx, &attribute.CreateAttributeInput{
		Name:       req.Name,
		Unit:       req.Unit,
		CategoryID: categoryID,
	})
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	resp := AttributeResponse{
		ID:         attr.ID,
		Name:       attr.Name,
		Unit:       attr.Unit,
		CategoryID: attr.CategoryID,
	}
	httputils.WriteJSON(w, http.StatusCreated, resp)
}

// ListAttributes godoc
// @Summary      Get all attributes for a category
// @Description  Получить все атрибуты в категории
// @Tags         categories
// @Param        categoryID path int true "Category ID"
// @Produce      json
// @Success      200 {array} AttributeResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
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

	resp := make(AttributeListResponse, len(attrs))
	for i, a := range attrs {
		resp[i] = AttributeResponse{
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
// @Success      200 {object} AttributeResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
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

	resp := AttributeResponse{
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
// @Param        data body UpdateAttributeReq true "Attribute update data"
// @Success      204 "No Content"
// @Failure      400 {object} dto.ErrorResponse
// @Failure      422 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
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

	var req UpdateAttributeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := h.validate.StructCtx(ctx, &req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			httputils.RespondValidationError(w, ve)
		} else {
			httputils.WriteError(w, http.StatusUnprocessableEntity, "invalid attribute data")
		}
		return
	}

	attr := &attr.Attribute{
		ID:         id,
		Name:       req.Name,
		Unit:       req.Unit,
		CategoryID: categoryID,
	}
	if err := h.srv.UpdateAttribute(ctx, attr); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteAttribute godoc
// @Summary      Delete an attribute
// @Description  Удалить атрибут по ID
// @Tags         categories
// @Security     BearerAuth
// @Param        categoryID path int true "Category ID"
// @Param        id path int true "Attribute ID"
// @Success      204 "No Content"
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
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
