package category

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/category/dto"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/interfaces"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, cat *catDom.Category) (*catDom.Category, error)
	GetCategory(ctx context.Context, id int64) (*catDom.Category, error)
	UpdateCategory(ctx context.Context, cat *catDom.Category) error
	DeleteCategory(ctx context.Context, id int64) error
	ListCategories(ctx context.Context) ([]catDom.Category, error)
}

type Handler struct {
	log *slog.Logger
	srv CategoryService
	val interfaces.Validator
}

type Deps struct {
	svc CategoryService
	val interfaces.Validator
	log *slog.Logger
}

func NewDeps(val interfaces.Validator, log *slog.Logger, svc CategoryService) (Deps, error) {
	if svc == nil {
		return Deps{}, errors.New("category handler: missing CategoryService")
	}
	if val == nil {
		return Deps{}, errors.New("category handler: missing validator")
	}
	if log == nil {
		return Deps{}, errors.New("category handler: missing logger")
	}
	return Deps{
		svc: svc,
		val: val,
		log: log.With("component", "restHTTP.category"),
	}, nil
}

func New(deps Deps) *Handler {
	return &Handler{
		srv: deps.svc,
		val: deps.val,
		log: deps.log,
	}
}

// CreateCategory godoc
//
//		@Summary		Create category
//		@Tags			categories
//		@Accept			json
//		@Produce		json
//	 @Security       BearerAuth
//		@Param			category	body		dto.CategoryCreateRequest	true	"Category to create"
//		@Success		201	 "Created"
//		@Failure		400			{object}	httputils.ErrorResponse
//		@Failure		409			{object}	httputils.ErrorResponse
//		@Failure		500			{object}	httputils.ErrorResponse
//		@Router			/api/admin/category [post]
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "category.Create")

	var input dto.CategoryCreateRequest
	defer func() {
		if cerr := r.Body.Close(); cerr != nil {
			log.Warn("body close failed", slog.Any("error", cerr))
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Warn("invalid JSON", slog.Any("error", err))
		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := h.val.StructCtx(ctx, &input); err != nil {
		log.Warn("validation failed", slog.Any("error", err))
		httputils.WriteError(w, http.StatusUnprocessableEntity, "invalid product data")
		return
	}

	cat := input.ToDomainCreate()
	_, err := h.srv.CreateCategory(ctx, cat)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// -----------------------------------------------------------------------------
// GetCategory godoc
//
//	@Summary		Get category by ID
//	@Tags			categories
//	@Produce		json
//
// @Param			id	path	int	true	"Category ID"
// @Success		200	{object}	dto.CategoryResponse
// @Failure		400	{object}	httputils.ErrorResponse
// @Failure		404	{object}	httputils.ErrorResponse
// @Failure		500	{object}	httputils.ErrorResponse
// @Router			/api/category/{id} [get]
func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	cat, err := h.srv.GetCategory(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	resp := dto.ToDTOResponse(cat)
	httputils.WriteJSON(w, http.StatusOK, resp)
}

// -----------------------------------------------------------------------------
// ListCategories godoc
//
//	@Summary		List categories
//	@Tags			categories
//	@Produce		json
//	@Success		200	{array}		dto.CategoryResponse
//	@Failure		500	{object}	httputils.ErrorResponse
//	@Router			/api/category [get]
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cats, err := h.srv.ListCategories(ctx)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	resp := make([]dto.CategoryResponse, 0, len(cats))
	for _, c := range cats {
		resp = append(resp, dto.ToDTOResponse(&c))
	}
	httputils.WriteJSON(w, http.StatusOK, resp)
}

// -----------------------------------------------------------------------------
// UpdateCategory godoc
//
//		@Summary		Update category
//		@Tags			categories
//		@Accept			json
//		@Produce		json
//	 	@Security       BearerAuth
//		@Param			id			path	int							true	"Category ID"
//		@Param			category	body	dto.CategoryUpdateRequest	true	"New name"
//		@Success		204
//		@Failure		400	{object}	httputils.ErrorResponse
//		@Failure		404	{object}	httputils.ErrorResponse
//		@Failure		500	{object}	httputils.ErrorResponse
//		@Router			/api/admin/category/{id} [put]
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "UpdateCategory")
	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	var input dto.CategoryUpdateRequest
	defer func() {
		if cerr := r.Body.Close(); cerr != nil {
			log.Warn("body close failed", slog.Any("error", cerr))
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := h.val.StructCtx(ctx, &input); err != nil {
		httputils.WriteError(w, http.StatusUnprocessableEntity, "invalid category data")
		return
	}

	cat := input.ToDomainUpdate(id)

	if err := h.srv.UpdateCategory(ctx, cat); err != nil {
		h.handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// -----------------------------------------------------------------------------
// DeleteCategory godoc
//
//		@Summary		Delete category
//		@Tags			categories
//		@Produce		json
//	 	@Security       BearerAuth
//		@Param			id	path	int	true	"Category ID"
//		@Success		204
//		@Failure		400	{object}	httputils.ErrorResponse
//		@Failure		404	{object}	httputils.ErrorResponse
//		@Failure		500	{object}	httputils.ErrorResponse
//		@Router			/api/admin/category/{id} [delete]
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	if err := h.srv.DeleteCategory(ctx, id); err != nil {
		h.handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, catDom.ErrCategoryNotFound):
		httputils.WriteError(w, http.StatusNotFound, "category not found")

	case errors.Is(err, catDom.ErrCategoryNameEmpty):
		httputils.WriteError(w, http.StatusUnprocessableEntity, "category name cannot be empty")

	// Слишком длинное имя
	case errors.Is(err, catDom.ErrCategoryNameTooLong):
		httputils.WriteError(w, http.StatusUnprocessableEntity, "category name must be at most 255 characters")

	// Попытка удалить категорию, которая используется
	case errors.Is(err, catDom.ErrCategoryInUse):
		httputils.WriteError(w, http.StatusConflict, "category is in use and cannot be deleted")

	// Всё прочее — Internal Server Error
	default:
		h.log.Error("service error", slog.Any("error", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
