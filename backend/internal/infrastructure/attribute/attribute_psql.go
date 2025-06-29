package attribute

import (
	"context"
	"fmt"
	"log/slog"

	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	e "github.com/Neimess/zorkin-store-project/internal/infrastructure/error"
	der "github.com/Neimess/zorkin-store-project/pkg/app_error"
	"github.com/Neimess/zorkin-store-project/pkg/database"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Deps struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewDeps(db *sqlx.DB, log *slog.Logger) (Deps, error) {
	if db == nil {
		return Deps{}, fmt.Errorf("attribute repository: missing database connection")
	}
	if log == nil {
		return Deps{}, fmt.Errorf("attribute repository: missing logger")
	}
	return Deps{
		db:  db,
		log: log.With("component", "PGAttributeRepository"),
	}, nil
}

type PGAttributeRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewPGAttributeRepository(deps Deps) *PGAttributeRepository {
	return &PGAttributeRepository{db: deps.db, log: deps.log}
}

func (r *PGAttributeRepository) SaveBatch(ctx context.Context, attrs []attr.Attribute) error {
	const query = `
    INSERT INTO attributes (name, unit, category_id)
    SELECT * FROM UNNEST($1::text[], $2::text[], $3::bigint[])
    RETURNING attribute_id
    `

	names := make([]string, len(attrs))
	units := make([]*string, len(attrs))
	catIDs := make([]int64, len(attrs))
	for i, a := range attrs {
		names[i] = a.Name
		units[i] = a.Unit
		catIDs[i] = a.CategoryID
	}

	returnedIDs := make([]int64, 0, len(attrs))

	err := r.withQuery(ctx, query, func() error {
		return r.db.SelectContext(ctx, &returnedIDs,
			query,
			pq.Array(names),
			pq.Array(units),
			pq.Array(catIDs),
		)
	})
	if err != nil {
		return r.mapPostgreSQLError(err)
	}

	if len(returnedIDs) != len(attrs) {
		return fmt.Errorf("%w: получили %d id, ожидали %d", der.ErrInternal, len(returnedIDs), len(attrs))
	}

	for i, id := range returnedIDs {
		attrs[i].ID = id
	}

	return nil
}

func (r *PGAttributeRepository) Save(ctx context.Context, attr *attr.Attribute) error {
	const query = `
	INSERT INTO attributes (name, unit, category_id)
	VALUES ($1, $2, $3)
	RETURNING attribute_id
	`
	var id int64

	err := r.withQuery(ctx, query, func() error {
		return r.db.QueryRowContext(ctx, query, attr.Name, attr.Unit, attr.CategoryID).Scan(&id)
	})
	if err != nil {
		return r.mapPostgreSQLError(err)
	}
	attr.ID = id
	return nil
}

func (r *PGAttributeRepository) GetByID(ctx context.Context, id int64) (*attr.Attribute, error) {
	const query = `
	SELECT attribute_id, name, unit, category_id
	FROM attributes
	WHERE attribute_id = $1
	`

	var attr attr.Attribute

	var raw attributeDB
	err := r.withQuery(ctx, query, func() error {
		return r.db.GetContext(ctx, &raw, query, id)
	})
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}

	attr = *raw.toDomain()
	return &attr, nil
}

func (r *PGAttributeRepository) FindByCategory(ctx context.Context, categoryID int64) ([]attr.Attribute, error) {
	const query = `
	SELECT attribute_id, name, unit, category_id
	FROM attributes
	WHERE category_id = $1
	ORDER BY name
	`

	var raws []attributeDB
	err := r.withQuery(ctx, query, func() error {
		return r.db.SelectContext(ctx, &raws, query, categoryID)
	})
	if err != nil {
		return nil, r.mapPostgreSQLError(err)
	}

	return rawListToDomain(raws), nil
}

func (r *PGAttributeRepository) Update(ctx context.Context, attr *attr.Attribute) error {
	const query = `
	UPDATE attributes
	SET name = $1, unit = $2, category_id = $3
	WHERE attribute_id = $4
	`

	err := r.withQuery(ctx, query, func() error {
		res, err := r.db.ExecContext(ctx, query, attr.Name, attr.Unit, attr.CategoryID, attr.ID)
		if err != nil {
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("%w: attribute with ID %d not found", der.ErrNotFound, attr.ID)
		}
		return nil
	})
	if err != nil {
		return r.mapPostgreSQLError(err)
	}

	return nil
}

func (r *PGAttributeRepository) Delete(ctx context.Context, id int64) error {
	const query = `
	DELETE FROM attributes
	WHERE attribute_id = $1
	`
	err := r.withQuery(ctx, query, func() error {
		res, err := r.db.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("%w: attribute with ID %d not found", der.ErrNotFound, id)
		}
		return nil
	})
	if err != nil {
		return r.mapPostgreSQLError(err)
	}
	return nil
}

func (r *PGAttributeRepository) withQuery(ctx context.Context, query string, fn func() error, extras ...slog.Attr) error {
	return database.WithQuery(ctx, r.log, query, fn, extras...)
}

func (r *PGAttributeRepository) mapPostgreSQLError(err error) error {
	return e.MapPostgreSQLError(r.log, err)
}
