package category

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/category"
)

type CategoryRepository interface {
	Create(ctx context.Context, cat *category.Category) (*category.Category, error)
	GetByID(ctx context.Context, id int64) (*category.Category, error)
	Update(ctx context.Context, id int64, newName string) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]category.Category, error)
}

type Service struct {
	repo CategoryRepository
	log  *slog.Logger
}

type Deps struct {
	repo CategoryRepository
}

func NewDeps(repo CategoryRepository) (*Deps, error) {
	if repo == nil {
		return nil, errors.New("category: missing CategoryRepository")
	}
	return &Deps{repo: repo}, nil
}

func New(d *Deps) *Service {
	return &Service{repo: d.repo, log: slog.Default().With("component", "service.category")}
}

func (s *Service) CreateCategory(ctx context.Context, cat *category.Category) (*category.Category, error) {
	if err := cat.Validate(); err != nil {
		return nil, err
	}

	cat, err := s.repo.Create(ctx, cat)
	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *Service) GetCategory(ctx context.Context, id int64) (*category.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, category.ErrCategoryNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("service: failed to retrieve category: %w", err)
	}
	return cat, nil
}

func (s *Service) UpdateCategory(ctx context.Context, cat *category.Category) error {
	if err := cat.Validate(); err != nil {
		return err
	}

	err := s.repo.Update(ctx, cat.ID, cat.Name)
	if err != nil {
		if errors.Is(err, category.ErrCategoryNotFound) {
			return err
		}
		return fmt.Errorf("service: failed to update category: %w", err)
	}
	return nil
}

func (s *Service) DeleteCategory(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, category.ErrCategoryNotFound) || errors.Is(err, category.ErrCategoryInUse) {
			return err
		}
		return fmt.Errorf("service: failed to delete category: %w", err)
	}
	return nil
}

func (s *Service) ListCategories(ctx context.Context) ([]category.Category, error) {
	cats, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: failed to list categories: %w", err)
	}
	return cats, nil
}
