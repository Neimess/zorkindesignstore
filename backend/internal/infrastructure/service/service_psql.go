package service

import (
	"context"
	"log/slog"

	domService "github.com/Neimess/zorkin-store-project/internal/domain/service"
	repoError "github.com/Neimess/zorkin-store-project/internal/infrastructure/error"
	"github.com/jmoiron/sqlx"
)

type PGServiceRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewPGServiceRepository(db *sqlx.DB, log *slog.Logger) *PGServiceRepository {
	if db == nil {
		panic("NewPGServiceRepository: db is nil")
	}
	return &PGServiceRepository{
		db:  db,
		log: log,
	}
}

func (r *PGServiceRepository) Create(ctx context.Context, s *domService.Service) (*domService.Service, error) {
	const q = `INSERT INTO services (name, description, price) VALUES ($1, $2, $3) RETURNING service_id`
	var id int64
	err := r.db.QueryRowContext(ctx, q, s.Name, s.Description, s.Price).Scan(&id)
	if err != nil {
		return nil, repoError.MapPostgreSQLError(r.log, err)
	}
	s.ID = id
	return s, nil
}

func (r *PGServiceRepository) Get(ctx context.Context, id int64) (*domService.Service, error) {
	const q = `SELECT service_id, name, description, price FROM services WHERE service_id = $1`
	var raw serviceDB
	err := r.db.GetContext(ctx, &raw, q, id)
	if err != nil {
		return nil, repoError.MapPostgreSQLError(r.log, err)
	}
	return raw.toDomain(), nil
}

func (r *PGServiceRepository) List(ctx context.Context) ([]domService.Service, error) {
	const q = `SELECT service_id, name, description, price FROM services`
	var raws []serviceDB
	err := r.db.SelectContext(ctx, &raws, q)
	if err != nil {
		return nil, repoError.MapPostgreSQLError(r.log, err)
	}
	return rawServiceListToDomain(raws), nil
}

func (r *PGServiceRepository) Update(ctx context.Context, s *domService.Service) (*domService.Service, error) {
	const q = `UPDATE services SET name = $1, description = $2, price = $3 WHERE service_id = $4`
	_, err := r.db.ExecContext(ctx, q, s.Name, s.Description, s.Price, s.ID)
	if err != nil {
		return nil, repoError.MapPostgreSQLError(r.log, err)
	}
	return s, nil
}

func (r *PGServiceRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM services WHERE service_id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return repoError.MapPostgreSQLError(r.log, err)
	}
	return nil
}

// --- Product <-> Service ---
func (r *PGServiceRepository) AddServicesToProduct(ctx context.Context, productID int64, serviceIDs []int64) error {
	if len(serviceIDs) == 0 {
		return nil
	}
	const q = `INSERT INTO product_services (product_id, service_id) SELECT * FROM UNNEST($1::bigint[], $2::bigint[]) ON CONFLICT DO NOTHING`
	prodIDs := make([]int64, len(serviceIDs))
	for i := range serviceIDs {
		prodIDs[i] = productID
	}
	_, err := r.db.ExecContext(ctx, q, prodIDs, serviceIDs)
	if err != nil {
		return repoError.MapPostgreSQLError(r.log, err)
	}
	return nil
}

func (r *PGServiceRepository) GetServicesByProduct(ctx context.Context, productID int64) ([]domService.Service, error) {
	const q = `SELECT s.service_id, s.name, s.description, s.price FROM services s JOIN product_services ps ON s.service_id = ps.service_id WHERE ps.product_id = $1`
	var raws []serviceDB
	err := r.db.SelectContext(ctx, &raws, q, productID)
	if err != nil {
		return nil, repoError.MapPostgreSQLError(r.log, err)
	}
	return rawServiceListToDomain(raws), nil
}
