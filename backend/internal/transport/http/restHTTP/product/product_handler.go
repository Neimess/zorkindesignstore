package product

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	"github.com/Neimess/zorkin-store-project/internal/domain/product"
	pSvc "github.com/Neimess/zorkin-store-project/internal/service/product"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/interfaces"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/product/dto"
	"github.com/Neimess/zorkin-store-project/pkg/httputils"
)

type ProductService interface {
	Create(ctx context.Context, product *product.Product) (int64, error)
	CreateWithAttrs(ctx context.Context, product *product.Product) (int64, error)
	GetDetailed(ctx context.Context, id int64) (*product.Product, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]product.Product, error)
}

type Deps struct {
	pSvc ProductService
	val  interfaces.Validator
}

func NewDeps(srv ProductService, val interfaces.Validator) (*Deps, error) {
	if srv == nil {
		return nil, errors.New("product handler: missing ProductService")
	}
	if val == nil {
		return nil, errors.New("product handler: missing validator")
	}
	return &Deps{
		pSvc: srv,
		val:  val,
	}, nil
}

type Handler struct {
	srv ProductService
	log *slog.Logger
	val interfaces.Validator
}

func New(d *Deps) *Handler {
	return &Handler{
		srv: d.pSvc,
		log: slog.Default().With("component", "transport.http.restHTTP.product"),
		val: d.val,
	}
}

// CreateProduct godoc
// @Summary      Create product
// @Description  Creates a new product and returns its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product  body      dto.ProductCreateRequest  true  "Product to create"
// @Success      201
// @Failure      400      {object}  httputils.ErrorResponse "Bad request"
// @Failure      401      {object}  httputils.ErrorResponse "Unauthorized access"
// @Failure      403      {object}  httputils.ErrorResponse "Forbidden access"
// @Failure      405	  {object}  httputils.ErrorResponse "Method not allowed, e.g. POST on GET endpoint"
// @Failure      404      {object}  httputils.ErrorResponse "Not found"
// @Failure      409      {object}  httputils.ErrorResponse "Conflict, e.g. duplicate product"
// @Failure      422      {object}  httputils.ErrorResponse "Unprocessable entity, e.g. validation error"
// @Failure      429      {object}  httputils.ErrorResponse "Too many requests, e.g. rate limiting"
// @Failure      500      {object}  httputils.ErrorResponse "Internal server error"
// @Router       /api/admin/product [post]
func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "transport.http.restHTTP.product.Create")
	defer func() {
		if cerr := r.Body.Close(); cerr != nil {
			log.Warn("body close failed", slog.Any("error", cerr))
		}
	}()

	var input dto.ProductCreateRequest
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
	product := mapCreateReqToDomain(&input)

	_, err := h.srv.Create(ctx, product)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// CreateWithAttributes godoc
// @Summary      Create product with attributes
// @Description  Creates a new product with its attributes and returns the created ID.
// @Tags         products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product  body      dto.ProductCreateRequest  true  "Product to create"
// @Success      201
// @Failure      400      {object}  httputils.ErrorResponse "Bad request"
// @Failure      401      {object}  httputils.ErrorResponse "Unauthorized access"
// @Failure      403      {object}  httputils.ErrorResponse "Forbidden access"
// @Failure      404      {object}  httputils.ErrorResponse "Not found"
// @Failure      405	  {object}  httputils.ErrorResponse "Method not allowed, e.g. POST on GET endpoint"
// @Failure      409      {object}  httputils.ErrorResponse "Conflict, e.g. duplicate product"
// @Failure      500      {object}  httputils.ErrorResponse "Internal server error"
// @Router       /api/admin/product/detailed [post]
func (h *Handler) CreateWithAttributes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "transport.http.restHTTP.product.GetDetailed")
	defer func() {
		if cerr := r.Body.Close(); cerr != nil {
			log.Warn("body close failed", slog.Any("error", cerr))
		}
	}()

	var req dto.ProductCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("invalid JSON", slog.Any("error", err))
		httputils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	domainProd := mapCreateReqToDomain(&req)

	_, err := h.srv.CreateWithAttrs(ctx, domainProd)
	if err != nil {
		switch {
		case errors.Is(err, pSvc.ErrInvalidAttribute):
			log.Warn("validation failed",
				slog.String("reason", "invalid attribute_id"),
			)
			httputils.WriteError(w, http.StatusBadRequest, "invalid product data")
			return
		default:
			log.Error("unexpected service error", slog.Any("error", err))
			httputils.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

// GetDetailedProduct godoc
// @Summary      Get product
// @Description  Returns a product with its attributes by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  dto.ProductResponse  "Returns product details"
// @Failure      400  {object}  httputils.ErrorResponse  "Invalid ID"
// @Failure      404  {object}  httputils.ErrorResponse  "Not found"
// @Failure      500  {object}  httputils.ErrorResponse  "Internal server error"
// @Router       /api/product/{id} [get]
func (h *Handler) GetDetailed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "transport.http.restHTTP.product.GetDetailed")
	id, err := httputils.IDFromURL(r, "id")

	if err != nil || id <= 0 {
		log.Warn("parse id error",
			slog.Any("product_id", id),
			slog.String("error", err.Error()))
		httputils.WriteError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	product, err := h.srv.GetDetailed(ctx, id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	resp := mapDomainToProductResponse(product)

	httputils.WriteJSON(w, http.StatusOK, resp)
}

// ListProductsByCategory godoc
// @Summary      List products by category
// @Description  Returns all products that belong to the specified category
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {array}   dto.ProductResponse  "List of products"
// @Failure      400  {object}  httputils.ErrorResponse    "Invalid ID"
// @Failure      401  {object}  httputils.ErrorResponse    "Unauthorized access"
// @Failure      403  {object}  httputils.ErrorResponse    "Forbidden access"
// @Failure      404  {object}  httputils.ErrorResponse    "Category not found"
// @Failure      405  {object}  httputils.ErrorResponse    "Method not allowed, e.g. POST on GET endpoint"
// @Failure      500  {object}  httputils.ErrorResponse    "Internal server error"
// @Router       /api/product/category/{id} [get]
func (h *Handler) ListByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.log.With("op", "transport.http.restHTTP.product.GetByCategoryID")
	categoryID, err := httputils.IDFromURL(r, "id")

	if err != nil || categoryID <= 0 {
		log.Warn("parse category_id error",
			slog.Any("category_id", categoryID),
			slog.String("error", err.Error()))
		httputils.WriteError(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	products, err := h.srv.GetByCategoryID(ctx, categoryID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	resp := make(dto.ProductListResponse, 0, len(products))
	for _, p := range products {
		resp = append(resp, *mapDomainToProductResponse(&p))
	}

	httputils.WriteJSON(w, http.StatusOK, resp)
}

func mapCreateReqToDomain(req *dto.ProductCreateRequest) *product.Product {
	p := &product.Product{
		Name:        req.Name,
		Price:       *req.Price,
		Description: req.Description,
		CategoryID:  *req.CategoryID,
		ImageURL:    req.ImageURL,
	}
	for _, a := range req.Attributes {
		p.Attributes = append(p.Attributes, product.ProductAttribute{
			AttributeID: *a.AttributeID,
			Value:       a.Value,
		})
	}
	return p
}

func mapDomainToProductResponse(p *product.Product) *dto.ProductResponse {
	resp := &dto.ProductResponse{
		ProductID:   p.ID,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
		CategoryID:  p.CategoryID,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt,
	}
	for _, pa := range p.Attributes {

		resp.Attributes = append(resp.Attributes, dto.ProductAttributeValueResponse{
			AttributeID: pa.AttributeID,
			Name:        pa.Attribute.Name,
			Unit:        pa.Attribute.Unit,
			Value:       pa.Value,
		})
	}
	return resp
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, pSvc.ErrBadCategoryID) || errors.Is(err, catDom.ErrCategoryNotFound):
		h.log.Warn("invalid category reference", slog.Any("error", err))
		httputils.WriteError(w, http.StatusBadRequest, "invalid or missing category")

	case errors.Is(err, pSvc.ErrInvalidAttribute):
		h.log.Warn("invalid attribute reference", slog.Any("error", err))
		httputils.WriteError(w, http.StatusUnprocessableEntity, "invalid attribute data")

	case errors.Is(err, pSvc.ErrProductNotFound):
		h.log.Warn("product not found", slog.Any("error", err))
		httputils.WriteError(w, http.StatusNotFound, "product not found")

	default:
		h.log.Error("unhandled service error", slog.Any("error", err))
		httputils.WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
