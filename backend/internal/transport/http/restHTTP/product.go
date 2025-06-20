package restHTTP

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/Neimess/zorkin-store-project/internal/service"
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
	log := ph.log.With("op", "transport.http.restHTTP.product.Create")
	defer r.Body.Close()

	var input domain.Product
	if err := easyjson.UnmarshalFromReader(r.Body, &input); err != nil {
		log.Warn("invalid JSON", slog.Any("error", err))
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := ph.srv.Create(ctx, &input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProductIsNil),
			errors.Is(err, service.ErrProductNameEmpty):
			log.Warn("validation failed", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		default:
			log.Error("unexpected service error", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	log.Info("product created", slog.Int64("product_id", id))
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
		log.Warn("parse id error",
			slog.Any("product_id", id),
			slog.String("error", err.Error()))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	product, err := ph.srv.GetDetailed(ctx, id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			log.Warn("product not found", slog.Any("product_id", id))
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
