package coefficients

import (
	"context"
	"log/slog"

	domCoeff "github.com/Neimess/zorkin-store-project/internal/domain/coefficients"
	repoError "github.com/Neimess/zorkin-store-project/internal/infrastructure/error"
	"github.com/jmoiron/sqlx"
)

type PGCoefficientsRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewPGCoefficientsRepository(db *sqlx.DB, log *slog.Logger) *PGCoefficientsRepository {
	if db == nil {
		panic("NewPGCoefficientsRepository: db is nil")
	}
	return &PGCoefficientsRepository{
		db:  db,
		log: log,
	}
}

func (r *PGCoefficientsRepository) Create(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error) {
	const q = `INSERT INTO coefficients (name, value) VALUES ($1, $2) RETURNING coefficient_id`
	var id int64
	err := r.withQuery(ctx, q, func() error {
		return r.db.QueryRowContext(ctx, q, c.Name, c.Value).Scan(&id)
	})
	if err != nil {
		return nil, repoError.MapPostgreSQLError(r.log, err)
	}
	c.ID = id
	return c, nil
}

func (r *PGCoefficientsRepository) Get(ctx context.Context, id int64) (*domCoeff.Coefficient, error) {
	const q = `SELECT coefficient_id, name, value FROM coefficients WHERE coefficient_id = $1`
	var raw coefficientDB
	err := r.withQuery(ctx, q, func() error {
		return r.db.GetContext(ctx, &raw, q, id)
	})
	if err != nil {
		return nil, repoError.MapPostgreSQLError(r.log, err)
	}
	return raw.toDomain(), nil
}

func (r *PGCoefficientsRepository) List(ctx context.Context) ([]domCoeff.Coefficient, error) {
	const q = `SELECT coefficient_id, name, value FROM coefficients`
	var raws []coefficientDB
	err := r.withQuery(ctx, q, func() error {
		return r.db.SelectContext(ctx, &raws, q)
	})
	if err != nil {
		return nil, repoError.MapPostgreSQLError(r.log, err)
	}
	return rawCoeffListToDomain(raws), nil
}

func (r *PGCoefficientsRepository) Update(ctx context.Context, c *domCoeff.Coefficient) (*domCoeff.Coefficient, error) {
	const q = `UPDATE coefficients SET name = $1, value = $2 WHERE coefficient_id = $3`
	err := r.withQuery(ctx, q, func() error {
		_, execErr := r.db.ExecContext(ctx, q, c.Name, c.Value, c.ID)
		return execErr
	})
	if err != nil {
		return nil, repoError.MapPostgreSQLError(r.log, err)
	}
	return c, nil
}

func (r *PGCoefficientsRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM coefficients WHERE coefficient_id = $1`
	err := r.withQuery(ctx, q, func() error {
		_, execErr := r.db.ExecContext(ctx, q, id)
		return execErr
	})
	if err != nil {
		return repoError.MapPostgreSQLError(r.log, err)
	}
	return nil
}

func (r *PGCoefficientsRepository) withQuery(ctx context.Context, query string, fn func() error, extras ...slog.Attr) error {
	r.log.Debug("query", slog.String("query", query))
	return fn()
}
