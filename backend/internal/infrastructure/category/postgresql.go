package category

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/Neimess/zorkin-store-project/internal/infrastructure/tx"
	"github.com/Neimess/zorkin-store-project/pkg/database"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type categoryDB struct {
	ID   int64  `db:"category_id"`
	Name string `db:"name"`
}

type categoryAttrRow struct {
	IsRequired   bool    `db:"is_required"`
	Priority     int     `db:"priority"`
	AttributeID  int64   `db:"attribute_id"`
	Name         string  `db:"name"`
	Slug         string  `db:"slug"`
	Unit         *string `db:"unit"`
	IsFilterable bool    `db:"is_filterable"`
}

const (
	foreignKeyViolationCode string = "23503"
	uniqueViolationCode     string = "23505"
)

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryExists        = errors.New("category already exists with this name")
	ErrCategoryInUse         = errors.New("category is in use and cannot be deleted")
	ErrCategoryDuplicateName = errors.New("category already exists with this name")
)

type CategoryRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	if db == nil {
		panic("NewCategoryRepository: db is nil")
	}
	return &CategoryRepository{
		db:  db,
		log: slog.Default().With("component", "CategoryRepository"),
	}
}

func (r *CategoryRepository) Create(
	ctx context.Context,
	c *domain.Category,
) (int64, error) {
	r.log.Debug("Create category", slog.String("op", "category.postgresql.Create"), slog.String("name", c.Name))
	return tx.RunInTx(ctx, r.db, func(tx *sqlx.Tx) (int64, error) {
		const insertCat = `
		  INSERT INTO categories (name)
		  VALUES ($1)
		  RETURNING category_id`
		if err := r.withQuery(ctx, insertCat, func() error {
			return tx.QueryRowContext(ctx, insertCat, c.Name).Scan(&c.ID)
		}); err != nil {
			var pgErr *pgconn.PgError
			switch {
			case errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode:
				return 0, ErrCategoryExists
			case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
				return 0, err
			default:
				return 0, fmt.Errorf("insert category failed: %w", err)
			}
		}

		for _, ca := range c.Attributes {
			var attrID int64

			const upsertAttr = `
			  INSERT INTO attributes (name, slug, unit, is_filterable)
			  VALUES ($1, $2, $3, $4)
			  ON CONFLICT (slug) DO UPDATE SET
			    name = EXCLUDED.name,
			    unit = EXCLUDED.unit,
			    is_filterable = EXCLUDED.is_filterable
			  RETURNING attribute_id`
			a := ca.Attr
			if err := tx.QueryRowContext(ctx, upsertAttr,
				a.Name, a.Slug, a.Unit, a.IsFilterable,
			).Scan(&attrID); err != nil {
				return 0, fmt.Errorf("upsert attribute %q: %w", a.Slug, err)
			}

			const insertLink = `
			  INSERT INTO category_attributes
			    (category_id, attribute_id, is_required, priority)
			  VALUES ($1, $2, $3, $4)`
			if _, err := tx.ExecContext(ctx,
				insertLink,
				c.ID, attrID,
				ca.IsRequired, ca.Priority,
			); err != nil {
				return 0, fmt.Errorf("insert category_attribute: %w", err)
			}
		}

		return c.ID, nil
	})
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	r.log.Debug("Get wit hid", slog.String("op", "category.postgresql.GetByID"), slog.String("id", fmt.Sprintf("%d", id)))

	const qCat = `
	  SELECT category_id, name
	    FROM categories
	   WHERE category_id = $1`
	var catRow categoryDB
	if err := r.withQuery(ctx, qCat, func() error {
		return r.db.GetContext(ctx, &catRow, qCat, id)
	}); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrCategoryNotFound
		case errors.Is(err, ctx.Err()):
			return nil, err
		default:
			return nil, fmt.Errorf("get category row: %w", err)
		}
	}

	const qAttr = `
	   SELECT
		   ca.is_required,
		   ca.priority,
		   a.attribute_id,
		   a.name,
		   a.slug,
		   a.unit,
		   a.is_filterable
		 FROM category_attributes ca
		 JOIN attributes a USING (attribute_id)
		WHERE ca.category_id = $1
		ORDER BY ca.priority`

	var rows []categoryAttrRow
	if err := r.withQuery(ctx, qAttr, func() error {
		return r.db.SelectContext(ctx, &rows, qAttr, id)
	}); err != nil {
		return nil, fmt.Errorf("get category attributes: %w", err)
	}

	cat := &domain.Category{
		ID:   catRow.ID,
		Name: catRow.Name,
	}
	cat.Attributes = mapToDomain(rows)
	return cat, nil
}

func (r *CategoryRepository) List(ctx context.Context) ([]domain.Category, error) {
	r.log.Debug("List categories", slog.String("op", "category.postgresql.List"))

	const q = `
SELECT
  c.category_id,
  c.name        AS category_name,
  ca.is_required,
  ca.priority,
  a.attribute_id,
  a.name        AS attr_name,
  a.slug,
  a.unit,
  a.is_filterable
FROM categories c
LEFT JOIN category_attributes ca ON ca.category_id = c.category_id
LEFT JOIN attributes a          ON a.attribute_id = ca.attribute_id
ORDER BY c.name, ca.priority;
`
	var rows *sqlx.Rows
	if err := r.withQuery(ctx, q, func() error {
		var err error
		rows, err = r.db.QueryxContext(ctx, q)
		return err
	}); err != nil {
		return nil, fmt.Errorf("list categories query failed: %w", err)
	}
	defer rows.Close()

	var result []domain.Category
	var current *domain.Category

	for rows.Next() {
		var row struct {
			CategoryID   int64          `db:"category_id"`
			CategoryName string         `db:"category_name"`
			IsRequired   sql.NullBool   `db:"is_required"`
			Priority     sql.NullInt64  `db:"priority"`
			AttributeID  sql.NullInt64  `db:"attribute_id"`
			AttrName     sql.NullString `db:"attr_name"`
			Slug         sql.NullString `db:"slug"`
			Unit         *string        `db:"unit"`
			IsFilterable sql.NullBool   `db:"is_filterable"`
		}
		if err := rows.StructScan(&row); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		if current == nil || current.ID != row.CategoryID {
			result = append(result, domain.Category{
				ID:   row.CategoryID,
				Name: row.CategoryName,
			})
			current = &result[len(result)-1]
		}

		if row.AttributeID.Valid {
			current.Attributes = append(current.Attributes, domain.CategoryAttribute{
				IsRequired: row.IsRequired.Bool,
				Priority:   int(row.Priority.Int64),
				Attr: domain.Attribute{
					ID:           row.AttributeID.Int64,
					Name:         row.AttrName.String,
					Slug:         row.Slug.String,
					Unit:         row.Unit,
					IsFilterable: row.IsFilterable.Bool,
				},
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}
	return result, nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *domain.Category) error {
	const qUpdate = `UPDATE categories SET name=$2 WHERE category_id=$1`
	var res sql.Result
	if err := r.withQuery(ctx, qUpdate, func() error {
		var err error
		res, err = r.db.ExecContext(ctx, qUpdate, c.ID, c.Name)
		return err
	}); err != nil {
		r.log.Debug("Update category failed", slog.String("op", "category.postgresql.Update"), slog.String("name", c.Name), slog.String("error", err.Error()))
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return fmt.Errorf("update category: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update category: failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrCategoryNotFound
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	const qDelete = `DELETE FROM categories WHERE category_id=$1`
	var res sql.Result
	if err := r.withQuery(ctx, qDelete, func() error {
		var err error
		res, err = r.db.ExecContext(ctx, qDelete, id)
		return err
	}); err != nil {
		r.log.Debug("Delete category failed", slog.String("op", "category.postgresql.Delete"), slog.Int64("category_id", id), slog.String("error", err.Error()))
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolationCode:
			return ErrCategoryInUse
		case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
			return err
		default:
			return fmt.Errorf("delete category: %w", err)
		}
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete category: failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrCategoryNotFound
	}
	return nil
}

func (r *CategoryRepository) withQuery(ctx context.Context, sql string, fn func() error) error {
	err := database.WithQuery(ctx, r.log, sql, fn)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return err
}

func mapToDomain(rows []categoryAttrRow) []domain.CategoryAttribute {
	out := make([]domain.CategoryAttribute, len(rows))
	for i, r := range rows {
		out[i] = domain.CategoryAttribute{
			IsRequired: r.IsRequired,
			Priority:   r.Priority,
			Attr: domain.Attribute{
				ID:           r.AttributeID,
				Name:         r.Name,
				Slug:         r.Slug,
				Unit:         r.Unit,
				IsFilterable: r.IsFilterable,
			},
		}
	}
	return out
}
