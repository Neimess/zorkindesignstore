package product

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	prodDom "github.com/Neimess/zorkin-store-project/internal/domain/product"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/product/dto"
	"github.com/Neimess/zorkin-store-project/pkg/http_utils"
)

type ProductService interface {
	Create(ctx context.Context, product *prodDom.Product) (*prodDom.Product, error)
	CreateWithAttrs(ctx context.Context, product *prodDom.Product) (*prodDom.Product, error)
	GetDetailed(ctx context.Context, id int64) (*prodDom.Product, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]prodDom.Product, error)
	Update(ctx context.Context, product *prodDom.Product) (*prodDom.Product, error)
	Delete(ctx context.Context, id int64) error
}

type Deps struct {
	pSvc ProductService
	log  *slog.Logger
}

func NewDeps(log *slog.Logger, srv ProductService) (Deps, error) {
	if srv == nil {
		return Deps{}, errors.New("product handler: missing ProductService")
	}
	if log == nil {
		return Deps{}, errors.New("product handler: missing logger")
	}
	return Deps{
		pSvc: srv,
		log:  log.With("component", "restHTTP.product"),
	}, nil
}

type Handler struct {
	srv ProductService
	log *slog.Logger
}

func New(d Deps) *Handler {
	return &Handler{
		srv: d.pSvc,
		log: d.log,
	}
}

// CreateProduct godoc
// @Summary Создать продукт вместе с аттрибутами
// @Description  Creates a new product with it's attributes and returns its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product  body      dto.ProductRequest  true  "Product to create (may contain services)"
// @Success      201	  {object}  dto.ProductResponse
// @Failure      400      {object}  http_utils.ErrorResponse "Bad request"
// @Failure      401      {object}  http_utils.ErrorResponse "Unauthorized access"
// @Failure      403      {object}  http_utils.ErrorResponse "Forbidden access"
// @Failure      405	  {object}  http_utils.ErrorResponse "Method not allowed, e.g. POST on GET endpoint"
// @Failure      404      {object}  http_utils.ErrorResponse "Not found"
// @Failure      409      {object}  http_utils.ErrorResponse "Conflict, e.g. duplicate product"
// @Failure      422      {object}  http_utils.ErrorResponse "Unprocessable entity, e.g. validation error"
// @Failure      429      {object}  http_utils.ErrorResponse "Too many requests, e.g. rate limiting"
// @Failure      500      {object}  http_utils.ErrorResponse "Internal server error"
// @Router       /api/admin/product [post]
func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "transport.http.restHTTP.product.Create")

	req, ok := http_utils.DecodeAndValidate[dto.ProductRequest](w, r, log)
	if !ok {
		return
	}
	domainProduct := req.MapCreateToDomain()

	product, err := h.srv.Create(ctx, domainProduct)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	resp := dto.MapDomainToProductResponse(product)
	w.Header().Set("Location", fmt.Sprintf("/api/product/%d", resp.ProductID))
	http_utils.WriteJSON(w, http.StatusCreated, resp)
}

// GetDetailedProduct godoc
// @Summary      Get product
// @Description  Returns a product with its attributes by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  dto.ProductResponse  "Returns product details"
// @Failure      400  {object}  http_utils.ErrorResponse  "Invalid ID"
// @Failure      404  {object}  http_utils.ErrorResponse  "Not found"
// @Failure      500  {object}  http_utils.ErrorResponse  "Internal server error"
// @Router       /api/product/{id} [get]
func (h *Handler) GetDetailed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "transport.http.restHTTP.product.GetDetailed")
	id, err := http_utils.IDFromURL(r, "id")

	if err != nil || id <= 0 {
		log.Warn("parse id error",
			slog.Any("product_id", id),
			slog.String("error", err.Error()))
		http_utils.WriteError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	product, err := h.srv.GetDetailed(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	resp := dto.MapDomainToProductResponse(product)
	http_utils.WriteJSON(w, http.StatusOK, resp)
}

// ListProductsByCategory godoc
// @Summary      List products by category
// @Description  Returns all products that belong to the specified category
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {array}   dto.ProductResponse  "List of products"
// @Failure      400  {object}  http_utils.ErrorResponse    "Invalid ID"
// @Failure      401  {object}  http_utils.ErrorResponse    "Unauthorized access"
// @Failure      403  {object}  http_utils.ErrorResponse    "Forbidden access"
// @Failure      404  {object}  http_utils.ErrorResponse    "Category not found"
// @Failure      405  {object}  http_utils.ErrorResponse    "Method not allowed, e.g. POST on GET endpoint"
// @Failure      500  {object}  http_utils.ErrorResponse    "Internal server error"
// @Router       /api/product/category/{id} [get]
func (h *Handler) ListByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "transport.http.restHTTP.product.GetByCategoryID")
	categoryID, err := http_utils.IDFromURL(r, "id")

	if err != nil || categoryID <= 0 {
		log.Warn("parse category_id error",
			slog.Any("category_id", categoryID),
			slog.String("error", err.Error()))
		http_utils.WriteError(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	products, err := h.srv.GetByCategoryID(ctx, categoryID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	resp := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		resp = append(resp, *dto.MapDomainToProductResponse(&p))
	}

	http_utils.WriteJSON(w, http.StatusOK, resp)
}

// UpdateProduct godoc
// @Summary Обновить продукт
// @Description  Update an existing product (and its attributes if provided)
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                         true  "Product ID"
// @Param        product  body      dto.ProductRequest    true  "Product data to update (may contain services)"
// @Success      200      {object}  dto.ProductResponse
// @Failure      400      {object}  http_utils.ErrorResponse   "Bad request"
// @Failure      404      {object}  http_utils.ErrorResponse   "Not found"
// @Failure      422      {object}  http_utils.ErrorResponse   "Validation error"
// @Failure      500      {object}  http_utils.ErrorResponse   "Internal server error"
// @Router       /api/admin/product/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "transport.http.restHTTP.product.Update")

	// 1) parse ID
	id, err := http_utils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		log.Warn("invalid product ID", slog.Any("id", id), slog.Any("error", err))
		http_utils.WriteError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	req, ok := http_utils.DecodeAndValidate[dto.ProductRequest](w, r, log)
	if !ok {
		return
	}

	p := req.MapUpdateToDomain(id)
	prodRes, err := h.srv.Update(ctx, p)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	resp := dto.MapDomainToProductResponse(prodRes)
	http_utils.WriteJSON(w, http.StatusOK, resp)
}

// Delete godoc
// @Summary      Delete a product
// @Description  Remove a product by ID
// @Tags         products
// @Security     BearerAuth
// @Param        id   path      int  true  "Product ID"
// @Success      204  "No Content"
// @Failure      400  {object}  http_utils.ErrorResponse  "Invalid ID"
// @Failure      404  {object}  http_utils.ErrorResponse  "Not found"
// @Failure      500  {object}  http_utils.ErrorResponse  "Internal server error"
// @Router       /api/admin/product/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "transport.http.restHTTP.product.Delete")

	id, err := http_utils.IDFromURL(r, "id")
	if err != nil || id <= 0 {
		log.Warn("invalid product ID", slog.Any("id", id), slog.Any("error", err))
		http_utils.WriteError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	if err := h.srv.Delete(ctx, id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, prodDom.ErrBadCategoryID):
		h.log.Warn("invalid category reference", slog.Any("error", err))
		http_utils.WriteError(w, http.StatusBadRequest, "invalid or missing category")
	case errors.Is(err, prodDom.ErrBadServiceID):
		h.log.Warn("invalid category reference", slog.Any("error", err))
		http_utils.WriteError(w, http.StatusBadRequest, "invalid service id")
	case errors.Is(err, prodDom.ErrInvalidAttribute):
		h.log.Warn("invalid attribute reference", slog.Any("error", err))
		http_utils.WriteError(w, http.StatusUnprocessableEntity, "invalid attribute data")

	case errors.Is(err, prodDom.ErrProductNotFound):
		h.log.Warn("product not found", slog.Any("error", err))
		http_utils.WriteError(w, http.StatusNotFound, "product not found")

	default:
		h.log.Error("unhandled service error", slog.Any("error", err))
		http_utils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
