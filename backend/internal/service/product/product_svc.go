package product

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/category"
	domProduct "github.com/Neimess/zorkin-store-project/internal/domain/product"
	utils "github.com/Neimess/zorkin-store-project/internal/utils/svc"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
)

type ProductRepository interface {
	Create(ctx context.Context, p *domProduct.Product) (*domProduct.Product, error)
	CreateWithAttrs(ctx context.Context, p *domProduct.Product) (*domProduct.Product, error)
	GetWithAttrs(ctx context.Context, id int64) (*domProduct.Product, error)
	ListByCategory(ctx context.Context, catID int64) ([]domProduct.Product, error)
	UpdateWithAttrs(ctx context.Context, p *domProduct.Product) (*domProduct.Product, error)
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo ProductRepository
	log  *slog.Logger
}

type Deps struct {
	repo ProductRepository
	log  *slog.Logger
}

func NewDeps(repo ProductRepository, log *slog.Logger) (*Deps, error) {
	if repo == nil {
		return nil, errors.New("product service: missing ProductRepository")
	}
	if log == nil {
		return nil, errors.New("product service: missing logger")
	}
	return &Deps{repo: repo, log: log.With("component", "service.product")}, nil
}
func New(d *Deps) *Service {
	return &Service{
		repo: d.repo,
		log:  d.log,
	}
}

func (s *Service) Create(ctx context.Context, p *domProduct.Product) (*domProduct.Product, error) {
	const op = "service.product.Create"
	log := s.log.With("op", op)

	prod, err := s.repo.Create(ctx, p)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound:              domProduct.ErrBadCategoryID,
			category.ErrCategoryNotFound: domProduct.ErrBadCategoryID,
			der.ErrValidation:            domProduct.ErrInvalidAttribute,
			der.ErrBadRequest:            domProduct.ErrInvalidAttribute,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}

	log.Info("product created", slog.Int64("product_id", p.ID))
	return prod, nil
}

// CreateWithAttrs создает продукт с атрибутами и услугами (batch-insert связей реализован в репозитории)
// TODO: если потребуется валидация услуг или бизнес-логика — добавить здесь
func (s *Service) CreateWithAttrs(ctx context.Context, p *domProduct.Product) (*domProduct.Product, error) {
	const op = "service.product.CreateWithAttrs"
	log := s.log.With("op", op)

	prod, err := s.repo.CreateWithAttrs(ctx, p)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound:              domProduct.ErrBadCategoryID,
			category.ErrCategoryNotFound: domProduct.ErrBadCategoryID,
			der.ErrValidation:            domProduct.ErrInvalidAttribute,
			der.ErrBadRequest:            domProduct.ErrInvalidAttribute,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}

	log.Info("product with attributes created", slog.Int64("product_id", prod.ID))
	return prod, nil
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

// Update обновляет продукт, его атрибуты и услуги (batch-insert связей реализован в репозитории)
// TODO: если потребуется валидация услуг или бизнес-логика — добавить здесь
func (s *Service) Update(ctx context.Context, p *domProduct.Product) (*domProduct.Product, error) {
	const op = "service.product.Update"
	log := s.log.With("op", op)

	prod, err := s.repo.UpdateWithAttrs(ctx, p)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound:               domProduct.ErrProductNotFound,
			domProduct.ErrProductNotFound: domProduct.ErrProductNotFound,
			der.ErrValidation:             domProduct.ErrInvalidAttribute,
			der.ErrBadRequest:             domProduct.ErrInvalidAttribute,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}

	log.Info("product updated with attributes", slog.Int64("product_id", p.ID))
	return prod, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	const op = "service.product.Delete"
	log := s.log.With("op", op)

	if err := s.repo.Delete(ctx, id); err != nil {
		mapping := map[error]error{
			der.ErrNotFound:               nil,
			domProduct.ErrProductNotFound: nil,
		}
		return utils.ErrorHandler(log, op, err, mapping)
	}

	log.Info("product deleted", slog.Int64("product_id", id))
	return nil
}
