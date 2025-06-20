package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
)

var (
	ErrCategoryAttributeRepoIsNil = errors.New("category attribute service: repo is nil")
)

type CategoryAttributeRepository interface {
	Create(ctx context.Context, ca *domain.CategoryAttribute) error
	GetByCategory(ctx context.Context, categoryID int64) ([]domain.CategoryAttribute, error)
	Update(ctx context.Context, ca *domain.CategoryAttribute) error
	Delete(ctx context.Context, categoryID, attributeID int64) error
}

type CategoryAttributeService struct {
	repo CategoryAttributeRepository
	log  *slog.Logger
}

func NewCategoryAttributeService(repo CategoryAttributeRepository, log *slog.Logger) *CategoryAttributeService {
	if repo == nil {
		panic(ErrCategoryAttributeRepoIsNil)
	}
	return &CategoryAttributeService{repo: repo, log: logger.WithComponent(log, "service.category_attribute")}
}

func (s *CategoryAttributeService) Create(ctx context.Context, ca *domain.CategoryAttribute) error {
	const op = "service.category_attribute.Create"
	log := s.log.With("op", op)

	if err := s.repo.Create(ctx, ca); err != nil {
		log.Error("failed to create category_attribute", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("category_attribute created", slog.Int64("category_id", ca.CategoryID), slog.Int64("attribute_id", ca.AttributeID))
	return nil
}

func (s *CategoryAttributeService) GetByCategory(ctx context.Context, categoryID int64) ([]domain.CategoryAttribute, error) {
	const op = "service.category_attribute.GetByCategory"
	log := s.log.With("op", op)

	attrs, err := s.repo.GetByCategory(ctx, categoryID)
	if err != nil {
		log.Error("failed to get category attributes", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("category_attributes fetched", slog.Int("count", len(attrs)), slog.Int64("category_id", categoryID))
	return attrs, nil
}

func (s *CategoryAttributeService) Update(ctx context.Context, ca *domain.CategoryAttribute) error {
	const op = "service.category_attribute.Update"
	log := s.log.With("op", op)

	if err := s.repo.Update(ctx, ca); err != nil {
		log.Error("failed to update category_attribute", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("category_attribute updated", slog.Int64("category_id", ca.CategoryID), slog.Int64("attribute_id", ca.AttributeID))
	return nil
}

func (s *CategoryAttributeService) Delete(ctx context.Context, categoryID, attributeID int64) error {
	const op = "service.category_attribute.Delete"
	log := s.log.With("op", op)

	if err := s.repo.Delete(ctx, categoryID, attributeID); err != nil {
		log.Error("failed to delete category_attribute", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("category_attribute deleted", slog.Int64("category_id", categoryID), slog.Int64("attribute_id", attributeID))
	return nil
}
