package category

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/Neimess/zorkin-store-project/internal/service"
	_ "github.com/Neimess/zorkin-store-project/internal/transport/dto"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
	"github.com/Neimess/zorkin-store-project/pkg/validator"

	"github.com/mailru/easyjson"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, name string, attrs []domain.CategoryAttribute) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
	List(ctx context.Context) ([]domain.Category, error)
	Update(ctx context.Context, c *domain.Category) error
	Delete(ctx context.Context, id int64) error
}

type CategoryHandler struct {
	log *slog.Logger
	srv CategoryService
}

func NewCategoryHandler(service CategoryService) *CategoryHandler {
	return &CategoryHandler{
		srv: service,
		log: slog.Default().With("component", "restHTTP.category"),
	}
}

var validate = validator.GetValidator()

// CreateCategory godoc
//
//		@Summary		Create category
//		@Tags			categories
//		@Accept			json
//		@Produce		json
//	 @Security       BearerAuth
//		@Param			category	body	CreateCategoryReq	true	"Category to create"
//		@Success		201
//		@Failure		400			{object}	dto.ErrorResponse
//		@Failure		409			{object}	dto.ErrorResponse
//		@Failure		500			{object}	dto.ErrorResponse
//		@Router			/api/admin/category [post]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "category.Create")
	defer r.Body.Close()

	var req CreateCategoryReq
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		log.Warn("invalid JSON", slog.Any("err", err))
		httputils.WriteError(w, 400, "invalid JSON")
		return
	}
	// if err := validate.StructCtx(ctx, &req); err != nil {
	// 	log.Warn("validation failed", slog.Any("err", err))
	// 	httputils.WriteError(w, 400, "validation failed")
	// 	return
	// }

	attrs := make([]domain.CategoryAttribute, len(req.Attributes))
	for i, a := range req.Attributes {
		attrs[i] = domain.CategoryAttribute{
			Attr: domain.Attribute{
				Name:         a.Name,
				Slug:         a.Slug,
				Unit:         a.Unit,
				IsFilterable: a.IsFilterable,
			},
			IsRequired: a.IsRequired,
			Priority:   a.Priority,
		}
	}
	_, err := h.srv.CreateCategory(ctx, req.Name, attrs)
	switch {
	case errors.Is(err, service.ErrCategoryExists):
		log.Warn("duplicate category", slog.String("name", req.Name))
		httputils.WriteError(w, 409, "category already exists")
		return
	case err != nil:
		log.Error("service error", slog.Any("err", err))
		httputils.WriteError(w, 500, "internal server error")
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
// @Success		200	{object}	CategoryResponse
// @Failure		400	{object}	dto.ErrorResponse
// @Failure		404	{object}	dto.ErrorResponse
// @Failure		500	{object}	dto.ErrorResponse
// @Router			/api/category/{id} [get]
func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "category.Get")

	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		log.Warn("bad id", slog.Any("err", err))
		httputils.WriteError(w, 400, "invalid id")
		return
	}

	cat, err := h.srv.GetByID(ctx, id)
	switch {
	case errors.Is(err, service.ErrCategoryNotFound):
		log.Warn("category not found", slog.Int64("id", id))
		httputils.WriteError(w, 404, "not found")
		return
	case err != nil:
		log.Error("service error", slog.Any("err", err))
		httputils.WriteError(w, 500, "internal server error")
		return
	}

	resp := CategoryResponse{ID: cat.ID, Name: cat.Name, Attributes: make([]AttributePayload, 0, len(cat.Attributes))}
	for _, attr := range cat.Attributes {
		resp.Attributes = append(resp.Attributes, AttributePayload{
			Name:         attr.Attr.Name,
			Slug:         attr.Attr.Slug,
			Unit:         attr.Attr.Unit,
			IsFilterable: attr.Attr.IsFilterable,
			IsRequired:   attr.IsRequired,
			Priority:     attr.Priority,
		})
	}
	httputils.WriteJSON(w, http.StatusOK, &resp)
}

// -----------------------------------------------------------------------------
// ListCategories godoc
//
//	@Summary		List categories
//	@Tags			categories
//	@Produce		json
//	@Success		200	{array}	CategoryResponseList
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/api/category [get]
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "category.List")

	cats, err := h.srv.List(ctx)
	if err != nil {
		log.Error("service error", slog.Any("err", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	var resp CategoryResponseList

	for _, c := range cats {
		attrs := make([]AttributePayload, 0, len(c.Attributes))
		for _, attr := range c.Attributes {
			attrs = append(attrs, AttributePayload{
				Name:         attr.Attr.Name,
				Slug:         attr.Attr.Slug,
				Unit:         attr.Attr.Unit,
				IsFilterable: attr.Attr.IsFilterable,
				IsRequired:   attr.IsRequired,
				Priority:     attr.Priority,
			})
		}
		resp.Categories = append(resp.Categories, CategoryResponse{ID: c.ID, Name: c.Name, Attributes: attrs})
	}
	resp.Total = int64(len(cats))
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
//		@Param			category	body	CategoryUpdateRequest	true	"New name"
//		@Success		204
//		@Failure		400	{object}	dto.ErrorResponse
//		@Failure		404	{object}	dto.ErrorResponse
//		@Failure		500	{object}	dto.ErrorResponse
//		@Router			/api/admin/category/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "category.Update")

	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		httputils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req CategoryUpdateRequest
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := validate.StructCtx(ctx, &req); err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "validation failed")
		return
	}

	if err := h.srv.Update(ctx, &domain.Category{ID: id, Name: req.Name}); err != nil {
		switch {
		case errors.Is(err, service.ErrCategoryNotFound):
			httputils.WriteError(w, http.StatusNotFound, "not found")
		default:
			log.Error("service error", slog.Any("err", err))
			httputils.WriteError(w, http.StatusInternalServerError, "internal")
		}
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
//		@Failure		400	{object}	dto.ErrorResponse
//		@Failure		404	{object}	dto.ErrorResponse
//		@Failure		500	{object}	dto.ErrorResponse
//		@Router			/api/admin/category/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "category.Delete")

	id, err := httputils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		httputils.WriteError(w, 400, "invalid id")
		return
	}
	if err := h.srv.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, service.ErrCategoryNotFound):
			httputils.WriteError(w, http.StatusNotFound, "not found")
		case errors.Is(err, service.ErrCategoryInUse):
			httputils.WriteError(w, http.StatusConflict, "category in use")
		default:
			log.Error("service error", slog.Any("err", err))
			httputils.WriteError(w, 500, "internal")
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
