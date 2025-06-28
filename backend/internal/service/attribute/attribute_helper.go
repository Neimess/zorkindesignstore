package attribute

import (
	"context"
	"errors"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	der "github.com/Neimess/zorkin-store-project/internal/domain/error"
)

func (s *Service) ensureCategory(ctx context.Context, categoryID int64) error {
	_, err := s.repoCat.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, der.ErrNotFound) {
			return domain.ErrCategoryNotFound
		}
		return err
	}
	return nil
}

func (s *Service) toDomain(in *CreateAttributeInput) *attr.Attribute {
	return &attr.Attribute{
		Name:       in.Name,
		Unit:       in.Unit,
		CategoryID: in.CategoryID,
	}
}
