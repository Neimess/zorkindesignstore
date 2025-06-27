package category

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	repoError "github.com/Neimess/zorkin-store-project/internal/infrastructure/error"
	"github.com/jmoiron/sqlx"
)

type categoryDB struct {
	ID   int64  `db:"category_id"`
	Name string `db:"name"`
}

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

func (r *PGCategoryRepository) Create(ctx context.Context, name string) (int64, error) {
	if name == "" {
		return 0, domain.ErrCategoryNameEmpty
	}
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO categories (name) VALUES ($1) RETURNING category_id`,
		name,
	).Scan(&id)
	if err != nil {
		return 0, repoError.MapPostgreSQLError(err)
	}
	return id, nil
}

func (r *PGCategoryRepository) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	var c domain.Category
	err := r.db.QueryRowContext(ctx,
		`SELECT category_id, name FROM categories WHERE category_id = $1`,
		id,
	).Scan(&c.ID, &c.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, repoError.MapPostgreSQLError(err)
	}
	return &c, nil
}

func (r *PGCategoryRepository) Update(ctx context.Context, id int64, newName string) error {
	if newName == "" {
		return domain.ErrCategoryNameEmpty
	}
	res, err := r.db.ExecContext(ctx,
		`UPDATE categories SET name = $1 WHERE category_id = $2`,
		newName, id,
	)
	if err != nil {
		return repoError.MapPostgreSQLError(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return repoError.MapPostgreSQLError(err)
	}
	if rows == 0 {
		return domain.ErrCategoryNotFound
	}
	return nil
}

func (r *PGCategoryRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM categories WHERE category_id = $1`,
		id,
	)
	if err != nil {
		return repoError.MapPostgreSQLError(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return repoError.MapPostgreSQLError(err)
	}
	if rows == 0 {
		return domain.ErrCategoryNotFound
	}
	return nil
}

func (r *PGCategoryRepository) List(ctx context.Context) ([]domain.Category, error) {
	var rows []categoryDB
	if err := r.db.SelectContext(ctx,
		&rows,
		`SELECT category_id, name FROM categories ORDER BY name`,
	); err != nil {
		return nil, repoError.MapPostgreSQLError(err)
	}

	cats := make([]domain.Category, len(rows))
	for i, c := range rows {
		cats[i] = domain.Category{ID: c.ID, Name: c.Name}
	}
	return cats, nil
}
