package product

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/category"
	domProduct "github.com/Neimess/zorkin-store-project/internal/domain/product"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrBadCategoryID    = errors.New("bad category ID")
	ErrInvalidAttribute = errors.New("invalid product attribute")
)

type ProductRepository interface {
	Create(ctx context.Context, p *domProduct.Product) (int64, error)
	CreateWithAttrs(ctx context.Context, p *domProduct.Product) (int64, error)
	GetWithAttrs(ctx context.Context, id int64) (*domProduct.Product, error)
	ListByCategory(ctx context.Context, catID int64) ([]domProduct.Product, error)
}

type Service struct {
	repo ProductRepository
	log  *slog.Logger
}

type Deps struct {
	repo ProductRepository
}

func NewDeps(repo ProductRepository) (*Deps, error) {
	if repo == nil {
		return nil, errors.New("product service: missing ProductRepository")
	}
	return &Deps{repo: repo}, nil
}
func New(d *Deps) *Service {
	return &Service{
		repo: d.repo,
		log:  slog.Default().With("component", "service.product"),
	}
}

func (s *Service) Create(ctx context.Context, p *domProduct.Product) (int64, error) {
	const op = "service.product.Create"
	log := s.log.With("op", op)

	id, err := s.repo.Create(ctx, p)
	if err != nil {
		switch {
		case errors.Is(err, der.ErrNotFound), errors.Is(err, category.ErrCategoryNotFound):
			log.Info("invalid category", slog.Int64("category_id", p.CategoryID))
			return 0, domProduct.ErrBadCategoryID

		case errors.Is(err, der.ErrValidation), errors.Is(err, der.ErrBadRequest):
			log.Info("invalid attribute data", slog.Any("error", err))
			return 0, domProduct.ErrInvalidAttribute

		default:
			log.Error("repo failed", slog.Any("error", err))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("product created", slog.Int64("product_id", id))
	return id, nil
}

func (s *Service) CreateWithAttrs(ctx context.Context, p *domProduct.Product) (int64, error) {
	const op = "service.product.CreateWithAttrs"
	log := s.log.With("op", op)

	id, err := s.repo.CreateWithAttrs(ctx, p)
	if err != nil {
		switch {
		case errors.Is(err, der.ErrNotFound), errors.Is(err, category.ErrCategoryNotFound):
			log.Warn("invalid category", slog.Int64("category_id", p.CategoryID))
			return 0, domProduct.ErrBadCategoryID

		case errors.Is(err, der.ErrValidation), errors.Is(err, der.ErrBadRequest):
			log.Warn("invalid attribute data", slog.Any("error", err))
			return 0, domProduct.ErrInvalidAttribute

		default:
			log.Error("repo failed", slog.Any("error", err))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("product with attributes created", slog.Int64("product_id", id))
	return id, nil
}

func (s *Service) GetDetailed(ctx context.Context, id int64) (*domProduct.Product, error) {
	const op = "service.product.GetDetailed"
	log := s.log.With("op", op)

	product, err := s.repo.GetWithAttrs(ctx, id)
	if err != nil {
		if errors.Is(err, der.ErrNotFound) {
			return nil, domProduct.ErrProductNotFound
		}
		log.Error("repo failed", slog.Int64("product_id", id), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("product retrieved", slog.Int64("product_id", product.ID), slog.String("name", product.Name))
	return product, nil
}

func (s *Service) GetByCategoryID(ctx context.Context, catID int64) ([]domProduct.Product, error) {
	const op = "service.product.GetByCategoryID"
	log := s.log.With("op", op)

	products, err := s.repo.ListByCategory(ctx, catID)
	if err != nil {
		if errors.Is(err, der.ErrNotFound) || errors.Is(err, category.ErrCategoryNotFound) {
			return nil, domProduct.ErrBadCategoryID
		}
		log.Error("repo failed", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("products retrieved", slog.Int64("category_id", catID), slog.Int("count", len(products)))
	return products, nil
}
