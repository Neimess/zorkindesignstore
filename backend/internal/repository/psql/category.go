package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category already exists with this name")
	ErrCategoryInUse    = errors.New("category is in use and cannot be deleted")
)

type CategoryRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewCategoryRepository(db *sqlx.DB, log *slog.Logger) *CategoryRepository {
	if db == nil {
		panic("NewCategoryRepository: db is nil")
	}
	return &CategoryRepository{
		db:  db,
		log: logger.WithComponent(log, "repo.category"),
	}
}

func (r *CategoryRepository) Create(ctx context.Context, c *domain.Category) (int64, error) {
	const q = `INSERT INTO categories(name) VALUES($1) RETURNING category_id`
	done := dbLog(ctx, r.log, q)

	var id int64
	err := r.db.QueryRowContext(ctx, q, c.Name).Scan(&id)
	done(err, slog.String("name", c.Name))

	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode:
		return 0, ErrCategoryExists
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		return 0, err
	case err != nil:
		return 0, fmt.Errorf("create category: %w", err)
	}
	c.ID = id
	return id, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	const q = `SELECT category_id, name FROM categories WHERE category_id=$1`
	done := dbLog(ctx, r.log, q)

	var c domain.Category
	err := r.db.GetContext(ctx, &c, q, id)
	done(err, slog.Int64("category_id", id))

	switch {
	case err == nil:
		return &c, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrCategoryNotFound
	case errors.Is(err, ctx.Err()):
		return nil, err
	default:
		return nil, fmt.Errorf("get category: %w", err)
	}
}

func (r *CategoryRepository) List(ctx context.Context) ([]domain.Category, error) {
	const q = `SELECT category_id, name FROM categories ORDER BY name`
	done := dbLog(ctx, r.log, q)

	var cats []domain.Category
	err := r.db.SelectContext(ctx, &cats, q)
	done(err)

	if err != nil {
		if errors.Is(err, ctx.Err()) {
			return nil, err
		}
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return cats, nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *domain.Category) error {
	const q = `UPDATE categories SET name=$2 WHERE category_id=$1`
	done := dbLog(ctx, r.log, q)

	res, err := r.db.ExecContext(ctx, q, c.ID, c.Name)
	done(err, slog.Int64("category_id", c.ID))

	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrCategoryNotFound
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM categories WHERE category_id=$1`
	done := dbLog(ctx, r.log, q)

	res, err := r.db.ExecContext(ctx, q, id)
	done(err, slog.Int64("category_id", id))
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolationCode:
		return ErrCategoryInUse
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		return err
	case err != nil:
		return fmt.Errorf("delete category: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrCategoryNotFound
	}
	return nil
}
