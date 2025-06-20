package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	repository "github.com/Neimess/zorkin-store-project/internal/repository/psql"
)

var (
	ErrCategoryIsNil     = errors.New("category is nil")
	ErrCategoryNameEmpty = errors.New("category name is empty")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrCategoryRepoIsNil = errors.New("category service: repo is nil")
	ErrCategoryExists	 = errors.New("category already exists with this name")
)

type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
	List(ctx context.Context) ([]domain.Category, error)
	Update(ctx context.Context, c *domain.Category) error
	Delete(ctx context.Context, id int64) error
}

type CategoryService struct {
	repo CategoryRepository
	log  *slog.Logger
}

func NewCategoryService(repo CategoryRepository, logger *slog.Logger) (*CategoryService, error) {
	if repo == nil {
		return nil, ErrCategoryRepoIsNil
	}
	if logger == nil {
		logger = silentLogger()
	}
	return &CategoryService{
		repo: repo,
		log:  logger,
	}, nil
}

func (cs *CategoryService) Create(ctx context.Context, c *domain.Category) (int64, error) {
	const op = "service.category.Create"
	log := cs.log.With("op", op)

	id, err := cs.repo.Create(ctx, c)
	switch {
	case errors.Is(err, repository.ErrDuplicateName):
		log.Info("category already exists with this name", slog.String("name", c.Name))
		return 0, ErrCategoryExists
	case err != nil:
		return 0, fmt.Errorf("%s: %w", op, err)
	default:
		log.Info("category created", slog.Int64("category_id", id), slog.String("name", c.Name))
	}
	return id, nil
}

func (cs *CategoryService) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	const op = "service.category.GetByID"
	log := cs.log.With("op", op)

	cat, err := cs.repo.GetByID(ctx, id)
	switch {
	case errors.Is(err, repository.ErrCategoryNotFound):
		return nil, ErrCategoryNotFound
	case errors.Is(err, context.Canceled):
		log.Info("context canceled while getting category", slog.Int64("category_id", id))
		return nil, context.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		log.Info("context deadline exceeded while getting category", slog.Int64("category_id", id))
		return nil, context.DeadlineExceeded
	case err != nil:
		return nil, fmt.Errorf("%s: %w", op, err)
	default:
		log.Info("product category retrieved",
			slog.Int64("category_id", cat.ID))
	}
	return cat, nil
}

func (cs *CategoryService) List(ctx context.Context) ([]domain.Category, error) {
	const op = "service.category.List"
	log := cs.log.With("op", op)

	cats, err := cs.repo.List(ctx)
	if err != nil {
		log.Error("failed to list categories", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("categories listed", slog.Int("count", len(cats)))
	return cats, nil
}

func (cs *CategoryService) Update(ctx context.Context, c *domain.Category) error {
	const op = "service.category.Update"
	log := cs.log.With("op", op)

	if err := cs.repo.Update(ctx, c); err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			log.Info("category not found for update", slog.Int64("category_id", c.ID))
			return ErrCategoryNotFound
		}
		log.Error("failed to update category", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("category updated", slog.Int64("category_id", c.ID), slog.String("name", c.Name))
	return nil
}

func (cs *CategoryService) Delete(ctx context.Context, id int64) error {
	const op = "service.category.Delete"
	log := cs.log.With("op", op)

	if err := cs.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			log.Info("category not found for delete", slog.Int64("category_id", id))
			return ErrCategoryNotFound
		}
		log.Error("failed to delete category", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("category deleted", slog.Int64("category_id", id))
	return nil
}
