package category

import (
	"context"
	"errors"
	"log/slog"

	catDom "github.com/Neimess/zorkin-store-project/internal/domain/category"
	repoError "github.com/Neimess/zorkin-store-project/internal/infrastructure/error"
	"github.com/Neimess/zorkin-store-project/pkg/app_error"
	"github.com/Neimess/zorkin-store-project/pkg/database"
	"github.com/jmoiron/sqlx"
)

type Deps struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewDeps(db *sqlx.DB, log *slog.Logger) (Deps, error) {
	if db == nil {
		return Deps{}, errors.New("category repository: missing database connection")
	}
	if log == nil {
		return Deps{}, errors.New("category repository: missing logger")
	}
	return Deps{
		db:  db,
		log: log.With("component", "PGCategoryRepository"),
	}, nil
}

type PGCategoryRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewPGCategoryRepository(d Deps) *PGCategoryRepository {
	return &PGCategoryRepository{
		db:  d.db,
		log: d.log,
	}
}

func (r *PGCategoryRepository) Create(ctx context.Context, cat *catDom.Category) (*catDom.Category, error) {
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

func (r *PGCategoryRepository) GetByID(ctx context.Context, id int64) (*catDom.Category, error) {
	var c catDom.Category
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

func (r *PGCategoryRepository) Update(ctx context.Context, id int64, newName string) (*catDom.Category, error) {
	const query = `UPDATE categories SET name = $1 WHERE category_id = $2 RETURNING category_id, name`
	var c catDom.Category
	err := r.withQuery(ctx, query, func() error {
		return r.db.QueryRowContext(ctx, query, newName, id).Scan(&c.ID, &c.Name)
	})
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	return &c, nil
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
		return app_error.ErrNotFound
	}
	return nil
}

func (r *PGCategoryRepository) List(ctx context.Context) ([]catDom.Category, error) {
	var rows []categoryDB
	const query = `SELECT category_id, name FROM categories ORDER BY name`
	if err := r.withQuery(ctx, query, func() error {
		return r.db.SelectContext(ctx, &rows, query)
	}); err != nil {
		return nil, r.mapPostgreSQLError(err)
	}

	cats := make([]catDom.Category, len(rows))
	for i, c := range rows {
		cats[i] = catDom.Category{ID: c.ID, Name: c.Name}
	}
	return cats, nil
}

func (r *PGCategoryRepository) withQuery(ctx context.Context, query string, fn func() error, extras ...slog.Attr) error {
	return database.WithQuery(ctx, r.log, query, fn, extras...)
}

func (r *PGCategoryRepository) mapPostgreSQLError(err error) error {
	return repoError.MapPostgreSQLError(r.log, err)
}
