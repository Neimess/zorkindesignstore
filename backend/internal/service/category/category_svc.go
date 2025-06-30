package category

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
)

type CategoryRepository interface {
	Create(ctx context.Context, cat *catDom.Category) (*catDom.Category, error)
	GetByID(ctx context.Context, id int64) (*catDom.Category, error)
	Update(ctx context.Context, id int64, newName string) (*catDom.Category, error)
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
		return nil, err
	}

	return cat, nil
}

func (s *Service) GetCategory(ctx context.Context, id int64) (*catDom.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, catDom.ErrCategoryNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("service: failed to retrieve category: %w", err)
	}
	return cat, nil
}

func (s *Service) UpdateCategory(ctx context.Context, cat *catDom.Category) (*catDom.Category, error) {
	if err := cat.Validate(); err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, cat.ID, cat.Name)
	if err != nil {
		if errors.Is(err, catDom.ErrCategoryNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("service: failed to update category: %w", err)
	}
	return updated, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, catDom.ErrCategoryNotFound) || errors.Is(err, catDom.ErrCategoryInUse) {
			return err
		}
		return fmt.Errorf("service: failed to delete category: %w", err)
	}
	return nil
}

func (s *Service) ListCategories(ctx context.Context) ([]catDom.Category, error) {
	cats, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: failed to list categories: %w", err)
	}
	return cats, nil
}
