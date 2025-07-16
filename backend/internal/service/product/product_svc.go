package product

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/category"
	domProduct "github.com/Neimess/zorkin-store-project/internal/domain/product"
	domService "github.com/Neimess/zorkin-store-project/internal/domain/service"
	utils "github.com/Neimess/zorkin-store-project/internal/utils/svc"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
)

type ProductRepository interface {
	Create(ctx context.Context, p *domProduct.Product) (*domProduct.Product, error)
	CreateWithAttrs(ctx context.Context, p *domProduct.Product) (*domProduct.Product, error)
	Get(ctx context.Context, id int64) (*domProduct.Product, error)
	ListByCategory(ctx context.Context, catID int64) ([]domProduct.Product, error)
	UpdateWithAttrs(ctx context.Context, p *domProduct.Product) (*domProduct.Product, error)
	Delete(ctx context.Context, id int64) error
}

type ServiceRepository interface {
	Get(ctx context.Context, id int64) (*domService.Service, error)
}

type Service struct {
	repoPrd ProductRepository
	repoSvc ServiceRepository
	log     *slog.Logger
}

type Deps struct {
	repoPrd ProductRepository
	repoSvc ServiceRepository
	log     *slog.Logger
}

func NewDeps(repoPrd ProductRepository, repoSvc ServiceRepository, log *slog.Logger) (*Deps, error) {
	if repoPrd == nil {
		return nil, errors.New("product service: missing ProductRepository")
	}
	if repoSvc == nil {
		return nil, errors.New("product service: missing ServiceRepository")
	}
	if log == nil {
		return nil, errors.New("product service: missing logger")
	}
	return &Deps{repoPrd: repoPrd, repoSvc: repoSvc, log: log.With("component", "service.product")}, nil
}
func New(d *Deps) *Service {
	return &Service{
		repoPrd: d.repoPrd,
		repoSvc: d.repoSvc,
		log:     d.log,
	}
}

func (s *Service) Create(ctx context.Context, p *domProduct.Product) (*domProduct.Product, error) {
	const op = "service.product.Create"
	log := s.log.With("op", op)

	prod, err := s.repoPrd.CreateWithAttrs(ctx, p)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound:               domProduct.ErrBadCategoryID,
			domService.ErrServiceNotFound: domProduct.ErrBadServiceID,
			category.ErrCategoryNotFound:  domProduct.ErrBadCategoryID,
			der.ErrValidation:             domProduct.ErrInvalidAttribute,
			der.ErrBadRequest:             domProduct.ErrInvalidAttribute,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}

	if err := s.fetchServices(ctx, prod); err != nil {
		mapping := map[error]error{
			der.ErrNotFound: domService.ErrServiceNotFound,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}
	log.Info("product created", slog.Int64("product_id", prod.ID))
	return prod, nil
}

// CreateWithAttrs создает продукт с атрибутами и услугами (batch-insert связей реализован в репозитории)
// TODO: если потребуется валидация услуг или бизнес-логика — добавить здесь
func (s *Service) CreateWithAttrs(ctx context.Context, p *domProduct.Product) (*domProduct.Product, error) {
	const op = "service.product.CreateWithAttrs"
	log := s.log.With("op", op)

	prod, err := s.repoPrd.CreateWithAttrs(ctx, p)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound:              domProduct.ErrBadCategoryID,
			category.ErrCategoryNotFound: domProduct.ErrBadCategoryID,
			der.ErrValidation:            domProduct.ErrInvalidAttribute,
			der.ErrBadRequest:            domProduct.ErrInvalidAttribute,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}

	if err := s.fetchServices(ctx, prod); err != nil {
		mapping := map[error]error{
			der.ErrNotFound: domService.ErrServiceNotFound,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}
	log.Info("product with attributes created", slog.Int64("product_id", prod.ID))
	return prod, nil
}

func (s *Service) GetDetailed(ctx context.Context, id int64) (*domProduct.Product, error) {
	const op = "service.product.GetDetailed"
	log := s.log.With("op", op)

	product, err := s.repoPrd.Get(ctx, id)
	if err != nil {
		if errors.Is(err, der.ErrNotFound) {
			return nil, domProduct.ErrProductNotFound
		}
		log.Error("repoPrd failed", slog.Int64("product_id", id), slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("product retrieved", slog.Int64("product_id", product.ID), slog.String("name", product.Name))
	return product, nil
}

func (s *Service) GetByCategoryID(ctx context.Context, catID int64) ([]domProduct.Product, error) {
	const op = "service.product.GetByCategoryID"
	log := s.log.With("op", op)

	products, err := s.repoPrd.ListByCategory(ctx, catID)
	if err != nil {
		if errors.Is(err, der.ErrNotFound) || errors.Is(err, category.ErrCategoryNotFound) {
			return nil, domProduct.ErrBadCategoryID
		}
		log.Error("repoPrd failed", slog.Any("error", err))
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

	prod, err := s.repoPrd.UpdateWithAttrs(ctx, p)
	if err != nil {
		mapping := map[error]error{
			der.ErrNotFound:               domProduct.ErrProductNotFound,
			domProduct.ErrProductNotFound: domProduct.ErrProductNotFound,
			der.ErrValidation:             domProduct.ErrInvalidAttribute,
			der.ErrBadRequest:             domProduct.ErrInvalidAttribute,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}

	if err := s.fetchServices(ctx, prod); err != nil {
		mapping := map[error]error{
			der.ErrNotFound: domService.ErrServiceNotFound,
		}
		return nil, utils.ErrorHandler(log, op, err, mapping)
	}
	log.Info("product updated with attributes", slog.Int64("product_id", p.ID))
	return prod, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	const op = "service.product.Delete"
	log := s.log.With("op", op)

	if err := s.repoPrd.Delete(ctx, id); err != nil {
		mapping := map[error]error{
			der.ErrNotFound: nil,
		}
		return utils.ErrorHandler(log, op, err, mapping)
	}

	log.Info("product deleted", slog.Int64("product_id", id))
	return nil
}

func (s *Service) fetchServices(ctx context.Context, p *domProduct.Product) error {
	var services []domService.Service
	for _, svc := range p.Services {
		service, err := s.repoSvc.Get(ctx, svc.ID)
		if err != nil {
			return err
		}
		services = append(services, *service)
	}
	p.Services = services
	return nil
}
