package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/jmoiron/sqlx"
)

var (
	ErrCategoryAttributeNotFound = errors.New("category attribute not found")
)

type CategoryAttributeRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewCategoryAttributeRepository(db *sqlx.DB, logger *slog.Logger) *CategoryAttributeRepository {
	return &CategoryAttributeRepository{
		db:  db,
		log: logger.With("component", "repo.category_attribute"),
	}
}

func (r *CategoryAttributeRepository) Create(ctx context.Context, ca *domain.CategoryAttribute) error {
	const q = `
		INSERT INTO category_attributes (category_id, attribute_id, is_required, priority)
		VALUES ($1, $2, $3, $4)
	`
	done := dbLog(ctx, r.log, q)

	_, err := r.db.ExecContext(ctx, q,
		ca.CategoryID,
		ca.AttributeID,
		ca.IsRequired,
		ca.Priority,
	)
	done(err,
		slog.Int64("category_id", ca.CategoryID),
		slog.Int64("attribute_id", ca.AttributeID),
	)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return fmt.Errorf("create category_attribute: %w", err)
	}
	return nil
}

func (r *CategoryAttributeRepository) GetByCategory(ctx context.Context, categoryID int64) ([]domain.CategoryAttribute, error) {
	const q = `
		SELECT category_id, attribute_id, is_required, priority
		FROM category_attributes
		WHERE category_id = $1
		ORDER BY priority
	`
	done := dbLog(ctx, r.log, q)

	var result []domain.CategoryAttribute
	err := r.db.SelectContext(ctx, &result, q, categoryID)
	done(err, slog.Int64("category_id", categoryID))

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, fmt.Errorf("get category_attributes by category: %w", err)
	}
	return result, nil
}

func (r *CategoryAttributeRepository) Update(ctx context.Context, ca *domain.CategoryAttribute) error {
	const q = `
		UPDATE category_attributes
		SET is_required = $3, priority = $4
		WHERE category_id = $1 AND attribute_id = $2
	`
	done := dbLog(ctx, r.log, q)

	res, err := r.db.ExecContext(ctx, q,
		ca.CategoryID,
		ca.AttributeID,
		ca.IsRequired,
		ca.Priority,
	)
	done(err,
		slog.Int64("category_id", ca.CategoryID),
		slog.Int64("attribute_id", ca.AttributeID),
	)
	if err != nil {
		return fmt.Errorf("update category_attribute: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrCategoryAttributeNotFound
	}
	return nil
}

func (r *CategoryAttributeRepository) Delete(ctx context.Context, categoryID, attributeID int64) error {
	const q = `
		DELETE FROM category_attributes
		WHERE category_id = $1 AND attribute_id = $2
	`
	done := dbLog(ctx, r.log, q)

	res, err := r.db.ExecContext(ctx, q, categoryID, attributeID)
	done(err,
		slog.Int64("category_id", categoryID),
		slog.Int64("attribute_id", attributeID),
	)
	if err != nil {
		return fmt.Errorf("delete category_attribute: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrCategoryAttributeNotFound
	}
	return nil
}
