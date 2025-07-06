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
	const query = `INSERT INTO categories (name, parent_id) VALUES ($1, $2) RETURNING category_id`
	// Prepare parent_id for null handling
	var parent interface{}
	if cat.ParentID != nil {
		parent = *cat.ParentID
	} else {
		parent = nil
	}
	err := r.db.QueryRowContext(ctx, query, cat.Name, parent).Scan(&id)
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	cat.ID = id
	return cat, nil
}

func (r *PGCategoryRepository) GetByID(ctx context.Context, id int64) (*catDom.Category, error) {
	var dbCat categoryDB
	const query = `SELECT category_id, name, parent_id FROM categories WHERE category_id = $1`
	err := r.db.GetContext(ctx, &dbCat, query, id)
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	cat := &catDom.Category{
		ID:   dbCat.ID,
		Name: dbCat.Name,
	}
	if dbCat.ParentID.Valid {
		pid := dbCat.ParentID.Int64
		cat.ParentID = &pid
	}
	return cat, nil
}

func (r *PGCategoryRepository) Update(ctx context.Context, cat *catDom.Category) (*catDom.Category, error) {
	const query = `UPDATE categories SET name = $1, parent_id = $2 WHERE category_id = $3 RETURNING category_id, name, parent_id`
	var dbCat categoryDB
	var parent interface{}
	if cat.ParentID != nil {
		parent = *cat.ParentID
	} else {
		parent = nil
	}
	err := r.db.GetContext(ctx, &dbCat, query, cat.Name, parent, cat.ID)
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	updated := &catDom.Category{ID: dbCat.ID, Name: dbCat.Name}
	if dbCat.ParentID.Valid {
		pid := dbCat.ParentID.Int64
		updated.ParentID = &pid
	}
	return updated, nil
}

func (r *PGCategoryRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE category_id = $1`, id)
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
	var dbCats []categoryDB
	const query = `SELECT category_id, name, parent_id FROM categories ORDER BY name`
	if err := r.db.SelectContext(ctx, &dbCats, query); err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	cats := make([]catDom.Category, len(dbCats))
	for i, dbCat := range dbCats {
		cats[i] = catDom.Category{ID: dbCat.ID, Name: dbCat.Name}
		if dbCat.ParentID.Valid {
			pid := dbCat.ParentID.Int64
			cats[i].ParentID = &pid
		}
	}
	return cats, nil
}

func (r *PGCategoryRepository) withQuery(ctx context.Context, query string, fn func() error, extras ...slog.Attr) error {
	return database.WithQuery(ctx, r.log, query, fn, extras...)
}

func (r *PGCategoryRepository) mapPostgreSQLError(err error) error {
	return repoError.MapPostgreSQLError(r.log, err)
}
