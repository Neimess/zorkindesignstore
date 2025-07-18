package attribute

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	attrDom "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	"github.com/Neimess/zorkin-store-project/internal/service/category"
	utils "github.com/Neimess/zorkin-store-project/internal/utils/svc"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
)

type AttributeRepository interface {
	SaveBatch(ctx context.Context, attrs []attrDom.Attribute) error
	Save(ctx context.Context, attr *attrDom.Attribute) error
	GetByID(ctx context.Context, id int64) (*attrDom.Attribute, error)
	FindByCategory(ctx context.Context, categoryID int64) ([]attrDom.Attribute, error)
	Update(ctx context.Context, attr *attrDom.Attribute) (*attrDom.Attribute, error)
	Delete(ctx context.Context, id int64) error
}

type Deps struct {
	repoAttr AttributeRepository
	repoCat  category.CategoryRepository
	log      *slog.Logger
}

func NewDeps(repoAttr AttributeRepository, repoCat category.CategoryRepository, log *slog.Logger) (*Deps, error) {
	if repoAttr == nil {
		return nil, errors.New("attribute service: missing AttributeRepository")
	}
	if repoCat == nil {
		return nil, errors.New("attribute service: missing CategoryRepository")
	}
	if log == nil {
		return nil, errors.New("attribute service: missing logger")
	}
	return &Deps{repoAttr: repoAttr, repoCat: repoCat, log: log.With("component", "service.attribute")}, nil
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
		log:      d.log,
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
		return utils.ErrorHandler(s.log, "service.attribute.CreateAttributesBatch", err, map[error]error{
			der.ErrConflict: attrDom.ErrAttributeAlreadyExists,
		})
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
		return nil, utils.ErrorHandler(s.log, "service.attribute.CreateAttribute", err, map[error]error{
			der.ErrConflict: attrDom.ErrAttributeAlreadyExists,
		})
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
		return nil, utils.ErrorHandler(s.log, "service.attribute.GetAttribute", err, map[error]error{
			der.ErrNotFound: attrDom.ErrAttributeNotFound,
		})
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
		return nil, utils.ErrorHandler(s.log, "service.attribute.ListAttributes", err, nil)
	}
	s.log.Info("ListAttributes succeeded", slog.Int("count", len(list)), slog.Int64("categoryID", categoryID))
	return list, nil
}

func (s *Service) UpdateAttribute(ctx context.Context, a *attrDom.Attribute) (*attrDom.Attribute, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}

	if err := s.ensureCategory(ctx, a.CategoryID); err != nil {
		return nil, err
	}

	updated, err := s.repoAttr.Update(ctx, a)
	if err != nil {
		s.log.Error("Update failed", slog.Any("error", err))
		return nil, utils.ErrorHandler(s.log, "service.attribute.UpdateAttribute", err, map[error]error{
			der.ErrNotFound: attrDom.ErrAttributeNotFound,
		})
	}
	s.log.Info("UpdateAttribute succeeded", slog.Int64("id", a.ID))
	return updated, nil
}

func (s *Service) DeleteAttribute(ctx context.Context, id int64) error {
	s.log.Debug("DeleteAttribute", slog.Int64("id", id))

	if err := s.repoAttr.Delete(ctx, id); err != nil {
		s.log.Error("Delete failed", slog.Any("error", err))
		return utils.ErrorHandler(s.log, "service.attribute.DeleteAttribute", err, map[error]error{
			der.ErrNotFound: nil,
		})
	}

	s.log.Info("DeleteAttribute succeeded", slog.Int64("id", id))
	return nil
}

func (s *Service) ensureCategory(ctx context.Context, categoryID int64) error {
	_, err := s.repoCat.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, der.ErrNotFound) {
			return catDom.ErrCategoryNotFound
		}
		s.log.Error("failed to fetch category", slog.Any("error", err))
		return fmt.Errorf("service: fetch category: %w", err)
	}
	return nil
}
