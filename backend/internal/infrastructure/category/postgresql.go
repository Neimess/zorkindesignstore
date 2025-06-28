package category

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/category"
	repoError "github.com/Neimess/zorkin-store-project/internal/infrastructure/error"
	"github.com/Neimess/zorkin-store-project/pkg/database"
	"github.com/jmoiron/sqlx"
)

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryExists        = errors.New("category already exists with this name")
	ErrCategoryInUse         = errors.New("category is in use and cannot be deleted")
	ErrCategoryDuplicateName = errors.New("category already exists with this name")
)

type PGCategoryRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewPGCategoryRepository(db *sqlx.DB) *PGCategoryRepository {
	if db == nil {
		panic("NewCategoryRepository: db is nil")
	}
	return &PGCategoryRepository{
		db:  db,
		log: slog.Default().With("component", "CategoryRepository"),
	}
}

func (r *PGCategoryRepository) Create(ctx context.Context, cat *category.Category) (*category.Category, error) {
	var id int64
	const query = `INSERT INTO categories (name) VALUES ($1) RETURNING category_id`
	err := r.withQuery(ctx, query, func() error {
		return r.db.QueryRowContext(ctx,
			query,
			cat.Name,
		).Scan(&id)
	})
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	cat.ID = id
	return cat, nil
}

func (r *PGCategoryRepository) GetByID(ctx context.Context, id int64) (*category.Category, error) {
	var c category.Category
	const query = `SELECT category_id, name FROM categories WHERE category_id = $1`
	err := r.withQuery(ctx, query, func() error {
		return r.db.QueryRowContext(ctx, query, id).
			Scan(&c.ID, &c.Name)
	})
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	return &c, nil
}

func (r *PGCategoryRepository) Update(ctx context.Context, id int64, newName string) error {
	const query = `UPDATE categories SET name = $1 WHERE category_id = $2`
	var res sql.Result
	err := r.withQuery(ctx, query, func() error {
		var err error
		res, err = r.db.ExecContext(ctx,
			query,
			newName, id,
		)
		return err
	})

	if err != nil {
		return r.mapPostgreSQLError(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return r.mapPostgreSQLError(err)
	}
	if rows == 0 {
		return category.ErrCategoryNotFound
	}
	return nil
}

func (r *PGCategoryRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM categories WHERE category_id = $1`,
		id,
	)
	if err != nil {
		return r.mapPostgreSQLError(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return r.mapPostgreSQLError(err)
	}
	if rows == 0 {
		return category.ErrCategoryNotFound
	}
	return nil
}

func (r *PGCategoryRepository) List(ctx context.Context) ([]category.Category, error) {
	var rows []categoryDB
	const query = `SELECT category_id, name FROM categories ORDER BY name`
	if err := r.withQuery(ctx, query, func() error {
		return r.db.SelectContext(ctx, &rows, query)
	}); err != nil {
		return nil, r.mapPostgreSQLError(err)
	}

	cats := make([]category.Category, len(rows))
	for i, c := range rows {
		cats[i] = category.Category{ID: c.ID, Name: c.Name}
	}
	return cats, nil
}

func (r *PGCategoryRepository) withQuery(ctx context.Context, query string, fn func() error, extras ...slog.Attr) error {
	return database.WithQuery(ctx, r.log, query, fn, extras...)
}

func (r *PGCategoryRepository) mapPostgreSQLError(err error) error {
	return repoError.MapPostgreSQLError(r.log, err)
}
