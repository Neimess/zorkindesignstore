package product

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	der "github.com/Neimess/zorkin-store-project/internal/domain/error"
	"github.com/Neimess/zorkin-store-project/internal/domain/product"
	repository "github.com/Neimess/zorkin-store-project/internal/infrastructure/product"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrBadCategoryID    = errors.New("bad category ID")
	ErrInvalidAttribute = errors.New("invalid product attribute")
)

type ProductRepository interface {
	Create(ctx context.Context, p *product.Product) (int64, error)
	CreateWithAttrs(ctx context.Context, p *product.Product) (int64, error)
	GetWithAttrs(ctx context.Context, id int64) (*product.Product, error)
	ListByCategory(ctx context.Context, catID int64) ([]product.Product, error)
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

func (s *Service) Create(ctx context.Context, product *product.Product) (int64, error) {
	const op = "service.product.Create"
	log := s.log.With("op", op)

	id, err := s.repo.Create(ctx, product)
	switch {
	case errors.Is(err, der.ErrNotFound), errors.Is(err, repository.ErrCategoryNotFound):
		log.Info("invalid category", slog.Int64("category_id", product.CategoryID), slog.Any("error", err))
		return 0, ErrBadCategoryID
	case errors.Is(err, der.ErrValidation), errors.Is(err, der.ErrBadRequest):
		log.Info("invalid attribute data", slog.Any("error", err))
		return 0, ErrInvalidAttribute
	case err != nil:
		log.Error("repo failed", slog.Int64("product_id", product.ID), slog.Any("error", err))
		return 0, fmt.Errorf("%s: %w", op, err)
	default:
		log.Info("product created", slog.Int64("product_id", id))
	}
	return id, nil
}

func (s *Service) CreateWithAttrs(ctx context.Context, product *product.Product) (int64, error) {
	const op = "service.product.CreateWithAttributes"
	log := s.log.With("op", op)

	id, err := s.repo.CreateWithAttrs(ctx, product)
	switch {
	case errors.Is(err, der.ErrNotFound), errors.Is(err, repository.ErrCategoryNotFound):
		log.Warn("invalid category", slog.Int64("category_id", product.CategoryID), slog.Any("error", err))
		return 0, ErrBadCategoryID
	case errors.Is(err, der.ErrValidation), errors.Is(err, der.ErrBadRequest):
		log.Warn("invalid attribute data", slog.Any("error", err))
		return 0, ErrInvalidAttribute
	case err != nil:
		log.Error("repo failed", slog.Int64("product_id", product.ID), slog.Any("error", err))
		return 0, fmt.Errorf("%s: %w", op, err)
	default:
		log.Info("product with attributes created", slog.Int64("product_id", id))
	}
	return id, nil
}

func (s *Service) GetDetailed(ctx context.Context, id int64) (*product.Product, error) {
	const op = "service.product.GetDetailed"
	log := s.log.With("op", op)

	product, err := s.repo.GetWithAttrs(ctx, id)
	switch {
	case errors.Is(err, repository.ErrProductNotFound), errors.Is(err, der.ErrNotFound):
		return nil, ErrProductNotFound
	case err != nil:
		log.Error("repo failed", slog.Int64("product_id", id), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("product retrieved", slog.Int64("product_id", product.ID), slog.String("name", product.Name))
	return product, nil
}

func (s *Service) GetByCategoryID(ctx context.Context, catID int64) ([]product.Product, error) {
	const op = "service.product.GetByCategoryID"
	log := s.log.With("op", op)

	products, err := s.repo.ListByCategory(ctx, catID)
	switch {
	case errors.Is(err, repository.ErrCategoryNotFound), errors.Is(err, der.ErrNotFound):
		return nil, ErrBadCategoryID
	case err != nil:
		log.Error("repo error", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	default:
		log.Info("products retrieved", slog.Int64("category_id", catID), slog.Int("count", len(products)))
	}
	return products, nil
}
