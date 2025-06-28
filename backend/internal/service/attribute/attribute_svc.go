package attribute

import (
	"context"
	"errors"
	"fmt"

	"log/slog"
	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	der "github.com/Neimess/zorkin-store-project/internal/domain/error"
	"github.com/Neimess/zorkin-store-project/internal/service/category"
	"github.com/Neimess/zorkin-store-project/pkg/validator"
)

type AttributeRepository interface {
	SaveBatch(ctx context.Context, attrs []*attr.Attribute) error
	Save(ctx context.Context, attr *attr.Attribute) error
	GetByID(ctx context.Context, id int64) (*attr.Attribute, error)
	FindByCategory(ctx context.Context, categoryID int64) ([]attr.Attribute, error)
	Update(ctx context.Context, attr *attr.Attribute) error
	Delete(ctx context.Context, id int64) error
}

type Deps struct {
	repoAttr AttributeRepository
	repoCat  category.CategoryRepository
}

func NewDeps(repoAttr AttributeRepository, repoCat category.CategoryRepository) (*Deps, error) {
	if repoAttr == nil {
		return nil, errors.New("attribute service: missing AttributeRepository")
	}
	if repoCat == nil {
		return nil, errors.New("attribute service: missing CategoryRepository")
	}
	return &Deps{
		repoAttr: repoAttr,
		repoCat:  repoCat,
	}, nil
}

type Service struct {
	repoAttr AttributeRepository
	repoCat  category.CategoryRepository
	log      *slog.Logger
}

func New(d *Deps) *Service {
	return &Service{
		repoAttr: d.repoAttr,
		repoCat:  d.repoCat,
		log:      slog.Default().With("component", "AttributeService")}
}

func (s *Service) CreateAttributesBatch(ctx context.Context, input CreateAttributesBatchInput) error {
	s.log.Debug("CreateAttributesBatch called", slog.Int("count", len(input)))
	if len(input) == 0 {
		return der.ErrBadRequest
	}
	if err := s.ensureCategory(ctx, input[0].CategoryID); err != nil {
		return err
	}
	attrs := make([]*attr.Attribute, len(input))
	for i, a := range input {
		attrs[i] = s.toDomain(a)
	}
	if err := s.repoAttr.SaveBatch(ctx, attrs); err != nil {
		s.log.Error("SaveBatch failed",
			slog.String("error", err.Error()))
		return fmt.Errorf("service: failed to batch save attributes: %w", err)
	}
	s.log.Info("CreateAttributesBatch succeeded", slog.Int("count", len(attrs)))
	return nil
}

func (s *Service) CreateAttribute(ctx context.Context, in *CreateAttributeInput) (*attr.Attribute, error) {
	s.log.Debug("CreateAttribute called", slog.String("name", in.Name), slog.Int64("categoryID", in.CategoryID))

	if err := s.ensureCategory(ctx, in.CategoryID); err != nil {
		return nil, err
	}

	validator := validator.GetValidator()
	err := validator.Struct(in)
	if err != nil {
		s.log.Error("validation failed",
			slog.String("input", in.Name), slog.String("error", err.Error()))
		return nil, fmt.Errorf("service: validation error: %w", err)
	}
	attr := s.toDomain(in)
	if err := s.repoAttr.Save(ctx, attr); err != nil {
		s.log.Error("Save failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("service: failed to save attribute: %w", err)
	}
	s.log.Info("CreateAttribute succeeded", slog.Int64("id", attr.ID))
	return attr, nil
}

func (s *Service) GetAttribute(ctx context.Context, categoryID, id int64) (*attr.Attribute, error) {
	s.log.Debug("GetAttribute", slog.Int64("categoryID", categoryID), slog.Int64("id", id))

	if err := s.ensureCategory(ctx, categoryID); err != nil {
		return nil, err
	}
	attr, err := s.repoAttr.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, der.ErrNotFound) {
			return nil, der.ErrNotFound
		}
		return nil, fmt.Errorf("get attribute: %w", err)
	}
	return attr, nil
}

func (s *Service) ListAttributes(ctx context.Context, categoryID int64) ([]attr.Attribute, error) {
	s.log.Debug("ListAttributes called", slog.Int64("categoryID", categoryID))
	if _, err := s.repoCat.GetByID(ctx, categoryID); err != nil {
		if errors.Is(err, der.ErrNotFound) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	attrs, err := s.repoAttr.FindByCategory(ctx, categoryID)
	if err != nil {
		s.log.Error("FindByCategory failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("service: failed to list attributes: %w", err)
	}
	s.log.Info("ListAttributes succeeded", slog.Int("count", len(attrs)), slog.Int64("categoryID", categoryID))
	return attrs, nil
}

func (s *Service) UpdateAttribute(ctx context.Context, attr *attr.Attribute) error {
	s.log.Debug("UpdateAttribute called", slog.Int64("id", attr.ID))
	if _, err := s.repoCat.GetByID(ctx, attr.CategoryID); err != nil {
		if errors.Is(err, der.ErrNotFound) {
			return domain.ErrCategoryNotFound
		}
		return fmt.Errorf("failed to get category: %w", err)
	}
	if err := s.repoAttr.Update(ctx, attr); err != nil {
		s.log.Error("Update failed", slog.String("error", err.Error()))
		return fmt.Errorf("service: failed to update attribute: %w", err)
	}
	s.log.Info("UpdateAttribute succeeded", slog.Int64("id", attr.ID))
	return nil
}

func (s *Service) DeleteAttribute(ctx context.Context, id int64) error {
	s.log.Debug("DeleteAttribute called", slog.Int64("id", id))
	if err := s.repoAttr.Delete(ctx, id); err != nil {
		s.log.Error("Delete failed", slog.String("error", err.Error()))
		return fmt.Errorf("service: failed to delete attribute: %w", err)
	}
	s.log.Info("DeleteAttribute succeeded", slog.Int64("id", id))
	return nil
}
