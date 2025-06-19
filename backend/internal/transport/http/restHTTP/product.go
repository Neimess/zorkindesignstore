package restHTTP

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	repository "github.com/Neimess/zorkin-store-project/internal/repository/psql"
	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
)

type ProductService interface {
	Create(ctx context.Context, product *domain.Product) (int64, error)
	GetDetailed(ctx context.Context, id int64) (*domain.Product, error)
	// GetByID(ctx context.Context, id int64) (*domain.ProductWithImages, error)
	// GetAll(ctx context.Context) ([]domain.ProductWithImages, error)
}

type ProductHandler struct {
	srv ProductService
	log *slog.Logger
}

func NewProductHandler(service ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{
		srv: service,
		log: logger.With("component", "product_http"),
	}
}

// CreateProduct godoc
// @Summary      Create product
// @Description  Creates a new product and returns its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      domain.Product  true  "Product to create"
// @Success      201      {object}  map[string]int64
// @Failure      400      {string}  string  "Bad request"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /api/product [post]
func (ph ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ph.log.With("op", "product.create")
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	log.Debug("RequestBody", "body", string(bodyBytes))
	if err != nil {
		log.Error("read body failed", "error", err)
		http.Error(w, "cannot read request", http.StatusBadRequest)
		return
	}
	var input domain.Product
	if err := easyjson.Unmarshal(bodyBytes, &input); err != nil {
		log.Warn("invalid json", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Debug("Product processed", slog.Any("input", input))
	id, err := ph.srv.Create(ctx, &input)
	if err != nil {
		var ve domain.ValidationError
		if errors.As(err, &ve) {
			log.Warn("validation error", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, repository.ErrCategoryNotExists) {
			log.Warn("invalid category", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Error("internal service error", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// GetDetailedProduct godoc
// @Summary      Get product
// @Description  Returns a product with its attributes by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  domain.Product
// @Failure      400  {string}  string  "Invalid ID"
// @Failure      404  {string}  string  "Not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/product/{id} [get]
func (ph *ProductHandler) GetDetailed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ph.log.With("op", "product.get_detailed")
	
	idStr := chi.URLParam(r, "id")
	log.Debug("parse id from param", "id", idStr)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	product, err := ph.srv.GetDetailed(ctx, id)
	log.Debug("Get product", "product", product)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			log.Warn("product not found", "product_id", id)
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		log.Error("internal service error", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Info("product retrieved", "product_id", product.ID)
	w.Header().Set("Content-Type", "application/json")
	_, _ = easyjson.MarshalToWriter(product, w)
}
