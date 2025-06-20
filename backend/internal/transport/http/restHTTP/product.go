package restHTTP

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/Neimess/zorkin-store-project/internal/service"
	"github.com/Neimess/zorkin-store-project/internal/transport/dto"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
	"github.com/Neimess/zorkin-store-project/pkg/validator"
	"github.com/mailru/easyjson"
)

type ProductService interface {
	Create(ctx context.Context, product *domain.Product) (int64, error)
	CreateWithAttrs(ctx context.Context, product *domain.Product) (int64, error)
	GetDetailed(ctx context.Context, id int64) (*domain.Product, error)
	// GetByCategoryID(ctx context.Context, categoryID int64) ([]domain.Product, error)
}

type ProductHandler struct {
	srv ProductService
	log *slog.Logger
}

func NewProductHandler(service ProductService, log *slog.Logger) *ProductHandler {
	if service == nil {
		panic("product handler: service is nil")
	}
	return &ProductHandler{
		srv: service,
		log: logger.WithComponent(log, "transport.http.restHTTP.product"),
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
// @Success      201      {object}  dto.IDResponse "Returns created product ID"
// @Failure      400      {object}  dto.ErrorResponse "Bad request"
// @Failure      401      {object}  dto.ErrorResponse "Unauthorized access"
// @Failure      403      {object}  dto.ErrorResponse "Forbidden access"
// @Failure      405	  {object}  dto.ErrorResponse "Method not allowed, e.g. POST on GET endpoint"
// @Failure      404      {object}  dto.ErrorResponse "Not found"
// @Failure      409      {object}  dto.ErrorResponse "Conflict, e.g. duplicate product"
// @Failure      422      {object}  dto.ErrorResponse "Unprocessable entity, e.g. validation error"
// @Failure      429      {object}  dto.ErrorResponse "Too many requests, e.g. rate limiting"
// @Failure      500      {object}  dto.ErrorResponse "Internal server error"
// @Router       /api/admin/product [post]
func (ph ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ph.log.With("op", "transport.http.restHTTP.product.Create")
	log.Debug("ctx", slog.Any("ctx", ctx))
	validator := validator.GetValidator()
	defer r.Body.Close()

	var input dto.ProductCreateRequest
	if err := easyjson.UnmarshalFromReader(r.Body, &input); err != nil {
		log.Warn("invalid JSON", slog.Any("error", err))
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	validator.StructCtx(ctx, &input)
	if err := validator.StructCtx(ctx, &input); err != nil {
		log.Warn("validation failed", slog.Any("error", err))
		writeError(w, http.StatusUnprocessableEntity, "invalid product data")
		return
	}
	product := mapCreateReqToDomain(&input)

	_, err := ph.srv.Create(ctx, product)
	switch {
	case errors.Is(err, service.ErrBadCategoryID):
		log.Warn("invalid foreign key in product",
			slog.Int64("category_id", product.CategoryID),
			slog.String("name", product.Name),
		)
		writeError(w, http.StatusBadRequest, "invalid foreign key in product")
		return
	case err != nil:
		log.Error("unexpected service error", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "internal server error")
	default:
		log.Error("unexpected service error", slog.Any("error", err))
		writeError(w, http.StatusInternalServerError, "internal server error")
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
// @Security     AdminAuth
// @Param        product  body      dto.ProductCreateRequest  true  "Product to create"
// @Success      201      {object}  dto.IDResponse  "Returns created product ID"
// @Failure      400      {object}  dto.ErrorResponse "Bad request"
// @Failure      401      {object}  dto.ErrorResponse "Unauthorized access"
// @Failure      403      {object}  dto.ErrorResponse "Forbidden access"
// @Failure      404      {object}  dto.ErrorResponse "Not found"
// @Failure      405	  {object}  dto.ErrorResponse "Method not allowed, e.g. POST on GET endpoint"
// @Failure      409      {object}  dto.ErrorResponse "Conflict, e.g. duplicate product"
// @Failure      500      {object}  dto.ErrorResponse "Internal server error"
// @Router       /api/admin/product/detailed [post]
func (ph *ProductHandler) CreateWithAttributes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ph.log.With("op", "transport.http.restHTTP.product.GetDetailed")
	defer r.Body.Close()

	var req dto.ProductCreateRequest
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		log.Warn("invalid JSON", slog.Any("error", err))
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	domainProd := mapCreateReqToDomain(&req)

	_, err := ph.srv.CreateWithAttrs(ctx, domainProd)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAttribute):
			log.Warn("validation failed",
				slog.String("reason", "invalid attribute_id"),
			)
			_, _ = writeError(w, http.StatusBadRequest, "invalid product data")
			return
		default:
			log.Error("unexpected service error", slog.Any("error", err))
			_, _ = writeError(w, http.StatusInternalServerError, "internal server error")
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
// @Failure      400  {object}  dto.ErrorResponse  "Invalid ID"
// @Failure      404  {object}  dto.ErrorResponse  "Not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /api/product/{id} [get]
func (ph *ProductHandler) GetDetailed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ph.log.With("op", "transport.http.restHTTP.product.GetDetailed")
	id, err := idFromUrlToInt64(r)

	if err != nil || id <= 0 {
		log.Warn("parse id error",
			slog.Any("product_id", id),
			slog.String("error", err.Error()))
		writeError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	product, err := ph.srv.GetDetailed(ctx, id)
	switch {
	case errors.Is(err, service.ErrProductNotFound):
		log.Warn("product not found", slog.Int64("product_id", id))
		writeError(w, 404, "product not found")
		return
	case err != nil:
		log.Error("unhandled service error", slog.Any("err", err))
		writeError(w, 500, "internal server error")
		return
	}

	resp := mapDomainToProductResponse(product)

	writeJSON(w, http.StatusOK, resp)
}



func mapCreateReqToDomain(req *dto.ProductCreateRequest) *domain.Product {
	p := &domain.Product{
		Name:        req.Name,
		Price:       *req.Price,
		Description: req.Description,
		CategoryID:  *req.CategoryID,
		ImageURL:    req.ImageURL,
	}
	for _, a := range req.Attributes {
		p.Attributes = append(p.Attributes, domain.ProductAttribute{
			AttributeID: *a.AttributeID,
			Value:       a.Value,
		})
	}
	return p
}

func mapDomainToProductResponse(p *domain.Product) *dto.ProductResponse {
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
			AttributeID:  pa.AttributeID,
			Name:         pa.Attribute.Name,
			Slug:         pa.Attribute.Slug,
			Unit:         pa.Attribute.Unit,
			IsFilterable: pa.Attribute.IsFilterable,
			Value:        pa.Value,
		})
	}
	return resp
}
