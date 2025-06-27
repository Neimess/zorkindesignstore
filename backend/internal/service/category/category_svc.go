package category

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, name string) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
	Update(ctx context.Context, id int64, newName string) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]domain.Category, error)
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

func (s *Service) CreateCategory(ctx context.Context, name string) (*domain.Category, error) {
	if name == "" {
		return nil, domain.ErrCategoryNameEmpty
	}

	id, err := s.repo.Create(ctx, name)
	if err != nil {
		return nil, err
	}

	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to fetch created category: %w", err)
	}
	return cat, nil
}

func (s *Service) GetCategory(ctx context.Context, id int64) (*domain.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("service: failed to retrieve category: %w", err)
	}
	return cat, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id int64, newName string) error {
	if newName == "" {
		return domain.ErrCategoryNameEmpty
	}

	err := s.repo.Update(ctx, id, newName)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			return err
		}
		return fmt.Errorf("service: failed to update category: %w", err)
	}
	return nil
}

func (s *Service) DeleteCategory(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) || errors.Is(err, domain.ErrCategoryInUse) {
			return err
		}
		return fmt.Errorf("service: failed to delete category: %w", err)
	}
	return nil
}

func (s *Service) ListCategories(ctx context.Context) ([]domain.Category, error) {
	cats, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: failed to list categories: %w", err)
	}
	return cats, nil
}
