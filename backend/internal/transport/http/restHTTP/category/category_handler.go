package category

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/category/dto"
	http_utils "github.com/Neimess/zorkin-store-project/pkg/http_utils"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, cat *catDom.Category) (*catDom.Category, error)
	GetCategory(ctx context.Context, id int64) (*catDom.Category, error)
	UpdateCategory(ctx context.Context, cat *catDom.Category) (*catDom.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	ListCategories(ctx context.Context) ([]catDom.Category, error)
}

type Handler struct {
	log *slog.Logger
	srv CategoryService
}

type Deps struct {
	svc CategoryService
	log *slog.Logger
}

func NewDeps(log *slog.Logger, svc CategoryService) (Deps, error) {
	if svc == nil {
		return Deps{}, errors.New("category handler: missing CategoryService")
	}
	if log == nil {
		return Deps{}, errors.New("category handler: missing logger")
	}
	return Deps{
		svc: svc,
		log: log.With("component", "restHTTP.category"),
	}, nil
}

func New(deps Deps) *Handler {
	return &Handler{
		srv: deps.svc,
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
//		@Param			category	body		dto.CategoryRequest	true	"Category to create"
//		@Success		201	 		{object}	dto.CategoryResponse
//		@Failure		400			{object}	http_utils.ErrorResponse
//		@Failure		409			{object}	http_utils.ErrorResponse
//		@Failure		500			{object}	http_utils.ErrorResponse
//		@Router			/api/admin/category [post]
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "category.Create")

	req, ok := http_utils.DecodeAndValidate[dto.CategoryRequest](w, r, log)
	if !ok {
		return
	}

	cat := req.ToDomainCreate()
	created, err := h.srv.CreateCategory(ctx, cat)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/api/category/%d", created.ID))
	http_utils.WriteJSON(w, http.StatusCreated, dto.ToDTOResponse(created))
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
// @Failure		400	{object}	http_utils.ErrorResponse
// @Failure		404	{object}	http_utils.ErrorResponse
// @Failure		500	{object}	http_utils.ErrorResponse
// @Router			/api/category/{id} [get]
func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := http_utils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	cat, err := h.srv.GetCategory(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	resp := dto.ToDTOResponse(cat)
	http_utils.WriteJSON(w, http.StatusOK, resp)
}

// -----------------------------------------------------------------------------
// ListCategories godoc
//
//	@Summary		List categories
//	@Tags			categories
//	@Produce		json
//	@Success		200	{array}		dto.CategoryResponse
//	@Failure		500	{object}	http_utils.ErrorResponse
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
	http_utils.WriteJSON(w, http.StatusOK, resp)
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
//		@Param			category	body	dto.CategoryRequest	true	"New name"
//		@Success		200	{object}	dto.CategoryResponse
//		@Failure		400	{object}	http_utils.ErrorResponse
//		@Failure		404	{object}	http_utils.ErrorResponse
//		@Failure		500	{object}	http_utils.ErrorResponse
//		@Router			/api/admin/category/{id} [put]
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "UpdateCategory")
	id, err := http_utils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	req, ok := http_utils.DecodeAndValidate[dto.CategoryRequest](w, r, log)
	if !ok {
		return
	}

	cat := req.ToDomainUpdate(id)

	updated, err := h.srv.UpdateCategory(ctx, cat)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	http_utils.WriteJSON(w, http.StatusOK, dto.ToDTOResponse(updated))
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
//		@Failure		400	{object}	http_utils.ErrorResponse
//		@Failure		404	{object}	http_utils.ErrorResponse
//		@Failure		500	{object}	http_utils.ErrorResponse
//		@Router			/api/admin/category/{id} [delete]
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := http_utils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid category id")
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
		http_utils.WriteError(w, http.StatusNotFound, "category not found")

	case errors.Is(err, catDom.ErrCategoryNameEmpty):
		http_utils.WriteError(w, http.StatusUnprocessableEntity, "category name cannot be empty")

	// Слишком длинное имя
	case errors.Is(err, catDom.ErrCategoryNameTooLong):
		http_utils.WriteError(w, http.StatusUnprocessableEntity, "category name must be at most 255 characters")

	// Попытка удалить категорию, которая используется
	case errors.Is(err, catDom.ErrCategoryInUse):
		http_utils.WriteError(w, http.StatusConflict, "category is in use and cannot be deleted")
	case errors.Is(err, catDom.ErrCategoryNameExists):
		http_utils.WriteError(w, http.StatusConflict, "category name already exist")
	// Всё прочее — Internal Server Error
	default:
		h.log.Error("service error", slog.Any("error", err))
		http_utils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
