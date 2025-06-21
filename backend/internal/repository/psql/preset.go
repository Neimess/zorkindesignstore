package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/Neimess/zorkin-store-project/pkg/database"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
)

// domain.ErrPresetNotFound declared in domain layer; keep local alias for convenience
var ErrPresetNotFound = errors.New("preset not found")

type PresetRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewPresetRepository(db *sqlx.DB, log *slog.Logger) *PresetRepository {
	if db == nil {
		panic("NewPresetRepository: db is nil")
	}
	return &PresetRepository{
		db:  db,
		log: logger.WithComponent(log, "repo.preset"),
	}
}

// Create inserts preset and its items in one transaction.
// Takes a fully‑formed *domain.Preset where Items contains product IDs.
func (r *PresetRepository) Create(ctx context.Context, p *domain.Preset) (int64, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if rec := recover(); rec != nil {
			_ = tx.Rollback()
			panic(rec)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	const insertPreset = `INSERT INTO presets (name, description, total_price, image_url)
                          VALUES ($1,$2,$3,$4)
                          RETURNING preset_id, created_at`

	var id int64
	var created time.Time
	err = database.WithQuery(ctx, r.log, insertPreset, func() error {
		return tx.QueryRowContext(ctx, insertPreset,
			p.Name, p.Description, p.TotalPrice, p.ImageURL,
		).Scan(&id, &created)
	})
	if err != nil {
		return 0, r.wrapError(err, "insert preset")
	}

	// Mass‑insert items if provided.
	if len(p.Items) > 0 {
		builder := sq.Insert("preset_items").
			Columns("preset_id", "product_id").
			PlaceholderFormat(sq.Dollar)
		for _, it := range p.Items {
			builder = builder.Values(id, it.ProductID)
		}
		sqlStr, args, buildErr := builder.ToSql()
		if buildErr != nil {
			return 0, fmt.Errorf("build insert items: %w", buildErr)
		}
		if err = database.WithQuery(ctx, r.log, sqlStr, func() error {
			_, execErr := tx.ExecContext(ctx, sqlStr, args...)
			return execErr
		}); err != nil {
			return 0, r.wrapError(err, "insert preset items")
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}

	p.ID, p.CreatedAt = id, created
	return id, nil
}

// Get returns preset with embedded Items slice.
func (r *PresetRepository) Get(ctx context.Context, id int64) (*domain.Preset, error) {
	const qPreset = `SELECT preset_id, name, description, total_price, image_url, created_at
                     FROM presets
                     WHERE preset_id = $1`

	var preset domain.Preset
	if err := database.WithQuery(ctx, r.log, qPreset, func() error {
		return r.db.GetContext(ctx, &preset, qPreset, id)
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPresetNotFound
		}
		return nil, fmt.Errorf("get preset: %w", err)
	}

	const qItems = `SELECT preset_item_id, preset_id, product_id
                    FROM preset_items WHERE preset_id = $1`
	if err := database.WithQuery(ctx, r.log, qItems, func() error {
		return r.db.SelectContext(ctx, &preset.Items, qItems, id)
	}); err != nil {
		return nil, fmt.Errorf("get preset items: %w", err)
	}
	return &preset, nil
}

// List fetches all presets with their items.
func (r *PresetRepository) ListDetailed(ctx context.Context) ([]domain.Preset, error) {
	const qPresets = `SELECT preset_id, name, description, total_price, image_url, created_at
                      FROM presets`
	var presets []domain.Preset
	if err := database.WithQuery(ctx, r.log, qPresets, func() error {
		return r.db.SelectContext(ctx, &presets, qPresets)
	}); err != nil {
		return nil, fmt.Errorf("list presets: %w", err)
	}
	if len(presets) == 0 {
		return presets, nil
	}

	ids := make([]int64, len(presets))
	for i, p := range presets {
		ids[i] = p.ID
	}

	const qItems = `
	SELECT
		pi.preset_item_id,
		pi.preset_id,
		pi.product_id,
		p.name  AS product_name,
		p.price AS product_price,
		p.image_url AS product_image_url
	FROM preset_items pi
	JOIN products p ON p.product_id = pi.product_id
	WHERE pi.preset_id = ANY($1)`

	type itemRow struct {
		PresetItemID int64          `db:"preset_item_id"`
		PresetID     int64          `db:"preset_id"`
		ProductID    int64          `db:"product_id"`
		Name         string         `db:"product_name"`
		Price        float64        `db:"product_price"`
		ImageURL     sql.NullString `db:"product_image_url"`
	}
	var rows []itemRow
	if err := database.WithQuery(ctx, r.log, qItems, func() error {
		return r.db.SelectContext(ctx, &rows, qItems, pq.Array(ids))
	}); err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}

	byID := make(map[int64]*domain.Preset, len(presets))
	for i := range presets {
		byID[presets[i].ID] = &presets[i]
	}
	for _, row := range rows {
		item := domain.PresetItem{
			ID:        row.PresetItemID,
			PresetID:  row.PresetID,
			ProductID: row.ProductID,
			Product: &domain.ProductSummary{
				ID:       row.ProductID,
				Name:     row.Name,
				Price:    row.Price,
				ImageURL: optionalString(row.ImageURL),
			},
		}
		byID[row.PresetID].Items = append(byID[row.PresetID].Items, item)
	}
	return presets, nil
}

func (r *PresetRepository) ListShort(ctx context.Context) ([]domain.Preset, error) {
	const q = `SELECT preset_id, name, description, total_price, image_url, created_at
			   FROM presets`
	var presets []domain.Preset
	if err := database.WithQuery(ctx, r.log, q, func() error {
		return r.db.SelectContext(ctx, &presets, q)
	}); err != nil {
		return nil, fmt.Errorf("list short presets: %w", err)
	}
	return presets, nil
}

// Delete removes preset and its items (ON DELETE CASCADE in DB or explicit delete).
func (r *PresetRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM presets WHERE preset_id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete preset: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrPresetNotFound
	}
	return nil
}

// Helper to translate pg errors to domain‑level ones.
func (r *PresetRepository) wrapError(err error, op string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case foreignKeyViolationCode:
			return fmt.Errorf("%s: foreign key", op)
		default:
			return fmt.Errorf("%s: %s", op, pgErr.Message)
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}
