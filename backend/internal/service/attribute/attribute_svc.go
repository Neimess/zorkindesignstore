package attribute

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	attrDom "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	"github.com/Neimess/zorkin-store-project/internal/service/category"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
)

type AttributeRepository interface {
	SaveBatch(ctx context.Context, attrs []attrDom.Attribute) error
	Save(ctx context.Context, attr *attrDom.Attribute) error
	GetByID(ctx context.Context, id int64) (*attrDom.Attribute, error)
	FindByCategory(ctx context.Context, categoryID int64) ([]attrDom.Attribute, error)
	Update(ctx context.Context, attr *attrDom.Attribute) error
	Delete(ctx context.Context, id int64) error
}

type Deps struct {
	RepoAttr AttributeRepository
	RepoCat  category.CategoryRepository
}

func NewDeps(repoAttr AttributeRepository, repoCat category.CategoryRepository) (*Deps, error) {
	if repoAttr == nil {
		return nil, errors.New("attribute service: missing AttributeRepository")
	}
	if repoCat == nil {
		return nil, errors.New("attribute service: missing CategoryRepository")
	}
	return &Deps{RepoAttr: repoAttr, RepoCat: repoCat}, nil
}

type Service struct {
	repoAttr AttributeRepository
	repoCat  category.CategoryRepository
	log      *slog.Logger
}

func New(d *Deps) *Service {
	return &Service{
		repoAttr: d.RepoAttr,
		repoCat:  d.RepoCat,
		log:      slog.Default().With("component", "AttributeService"),
	}
}

func (s *Service) CreateAttributesBatch(ctx context.Context, categoryID int64, attrs []attrDom.Attribute) error {
	if len(attrs) == 0 {
		return attrDom.ErrBatchEmpty
	}
	if err := s.ensureCategory(ctx, categoryID); err != nil {
		return err
	}

	if err := s.repoAttr.SaveBatch(ctx, attrs); err != nil {
		s.log.Error("SaveBatch failed", slog.Any("error", err))
		if errors.Is(err, der.ErrConflict) {
			return attrDom.ErrAttributeAlreadyExists
		}
		return fmt.Errorf("service: batch save attributes: %w", err)
	}

	s.log.Info("CreateAttributesBatch succeeded", slog.Int("count", len(attrs)))
	return nil
}

func (s *Service) CreateAttribute(ctx context.Context, categoryID int64, a *attrDom.Attribute) (*attrDom.Attribute, error) {
	if err := s.ensureCategory(ctx, categoryID); err != nil {
		return nil, err
	}

	if err := s.repoAttr.Save(ctx, a); err != nil {
		s.log.Error("Save failed", slog.Any("error", err))
		switch {
		case errors.Is(err, der.ErrConflict):
			return nil, attrDom.ErrAttributeAlreadyExists
		case errors.Is(err, der.ErrNotFound):
			return nil, catDom.ErrCategoryNotFound
		default:
			return nil, fmt.Errorf("service: save attribute: %w", err)
		}
	}

	s.log.Info("CreateAttribute succeeded", slog.Int64("id", a.ID))
	return a, nil
}

func (s *Service) GetAttribute(ctx context.Context, categoryID, id int64) (*attrDom.Attribute, error) {
	s.log.Debug("GetAttribute", slog.Int64("categoryID", categoryID), slog.Int64("id", id))

	if err := s.ensureCategory(ctx, categoryID); err != nil {
		return nil, err
	}

	a, err := s.repoAttr.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, der.ErrNotFound) {
			return nil, attrDom.ErrAttributeNotFound
		}
		return nil, fmt.Errorf("service: get attribute: %w", err)
	}
	return a, nil
}

func (s *Service) ListAttributes(ctx context.Context, categoryID int64) ([]attrDom.Attribute, error) {
	s.log.Debug("ListAttributes", slog.Int64("categoryID", categoryID))

	if err := s.ensureCategory(ctx, categoryID); err != nil {
		return nil, err
	}

	list, err := s.repoAttr.FindByCategory(ctx, categoryID)
	if err != nil {
		s.log.Error("FindByCategory failed", slog.Any("error", err))
		return nil, fmt.Errorf("service: list attributes: %w", err)
	}
	s.log.Info("ListAttributes succeeded", slog.Int("count", len(list)), slog.Int64("categoryID", categoryID))
	return list, nil
}

func (s *Service) UpdateAttribute(ctx context.Context, a *attrDom.Attribute) error {
	s.log.Debug("UpdateAttribute", slog.Int64("id", a.ID))

	if err := s.ensureCategory(ctx, a.CategoryID); err != nil {
		return err
	}

	if err := s.repoAttr.Update(ctx, a); err != nil {
		s.log.Error("Update failed", slog.Any("error", err))
		if errors.Is(err, der.ErrNotFound) {
			return attrDom.ErrAttributeNotFound
		}
		return fmt.Errorf("service: update attribute: %w", err)
	}
	s.log.Info("UpdateAttribute succeeded", slog.Int64("id", a.ID))
	return nil
}

func (s *Service) DeleteAttribute(ctx context.Context, id int64) error {
	s.log.Debug("DeleteAttribute", slog.Int64("id", id))

	if err := s.repoAttr.Delete(ctx, id); err != nil {
		s.log.Error("Delete failed", slog.Any("error", err))
		if errors.Is(err, der.ErrNotFound) {
			return attrDom.ErrAttributeNotFound
		}
		return fmt.Errorf("service: delete attribute: %w", err)
	}
	s.log.Info("DeleteAttribute succeeded", slog.Int64("id", id))
	return nil
}
