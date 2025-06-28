package attribute

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
)

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
