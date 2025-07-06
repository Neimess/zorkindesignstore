package category

import (
	"context"
	"errors"
	"log/slog"

	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	utils "github.com/Neimess/zorkin-store-project/internal/utils/svc"
	"github.com/Neimess/zorkin-store-project/pkg/app_error"
)

type CategoryRepository interface {
	Create(ctx context.Context, cat *catDom.Category) (*catDom.Category, error)
	GetByID(ctx context.Context, id int64) (*catDom.Category, error)
	Update(ctx context.Context, cat *catDom.Category) (*catDom.Category, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]catDom.Category, error)
}

type Service struct {
	repo CategoryRepository
	log  *slog.Logger
}

type Deps struct {
	repo CategoryRepository
	log  *slog.Logger
}

func NewDeps(repo CategoryRepository, log *slog.Logger) (*Deps, error) {
	if repo == nil {
		return nil, errors.New("category: missing CategoryRepository")
	}
	if log == nil {
		return nil, errors.New("category: missing logger")
	}
	return &Deps{repo: repo, log: log.With("component", "service.category")}, nil
}

func New(d *Deps) *Service {
	return &Service{repo: d.repo, log: d.log}
}

func (s *Service) CreateCategory(ctx context.Context, cat *catDom.Category) (*catDom.Category, error) {
	if err := cat.Validate(); err != nil {
		return nil, err
	}

	cat, err := s.repo.Create(ctx, cat)
	if err != nil {
		return nil, utils.ErrorHandler(s.log, "service.category.CreateCategory", err, map[error]error{
			app_error.ErrConflict: catDom.ErrCategoryNameExists,
		})
	}

	return cat, nil
}

func (s *Service) GetCategory(ctx context.Context, id int64) (*catDom.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, utils.ErrorHandler(s.log, "service.category.GetCategory", err, map[error]error{
			catDom.ErrCategoryNotFound: catDom.ErrCategoryNotFound,
		})
	}
	return cat, nil
}

func (s *Service) UpdateCategory(ctx context.Context, cat *catDom.Category) (*catDom.Category, error) {
	if err := cat.Validate(); err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, cat)
	if err != nil {
		return nil, utils.ErrorHandler(s.log, "service.category.UpdateCategory", err, map[error]error{
			catDom.ErrCategoryNotFound: catDom.ErrCategoryNotFound,
		})
	}
	return updated, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return utils.ErrorHandler(s.log, "service.category.DeleteCategory", err, map[error]error{
			catDom.ErrCategoryNotFound: nil,
			catDom.ErrCategoryInUse:    catDom.ErrCategoryInUse,
		})
	}
	return nil
}

func (s *Service) ListCategories(ctx context.Context) ([]catDom.Category, error) {
	cats, err := s.repo.List(ctx)
	if err != nil {
		return nil, utils.ErrorHandler(s.log, "service.category.ListCategories", err, nil)
	}
	return cats, nil
}
