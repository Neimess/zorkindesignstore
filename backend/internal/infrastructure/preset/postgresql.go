package preset

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/preset"
	"github.com/Neimess/zorkin-store-project/internal/domain/product"
	repoError "github.com/Neimess/zorkin-store-project/internal/infrastructure/error"
	"github.com/Neimess/zorkin-store-project/pkg/database"
	"github.com/Neimess/zorkin-store-project/pkg/database/tx"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
)

type PGPresetRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewPGPresetRepository(db *sqlx.DB, log *slog.Logger) *PGPresetRepository {
	if db == nil {
		panic("NewPresetRepository: db is nil")
	}
	return &PGPresetRepository{
		db:  db,
		log: logger.WithComponent(log, "repo.preset"),
	}
}

func (r *PGPresetRepository) Create(ctx context.Context, p *preset.Preset) (int64, error) {
	r.log.Debug("Create preset",
		slog.String("op", "preset.postgresql.Create"),
		slog.String("preset_name", p.Name),
	)

	return tx.RunInTx(ctx, r.db, func(tx *sqlx.Tx) (int64, error) {
		const insertPreset = `
			INSERT INTO presets (name, description, total_price, image_url)
			VALUES ($1, $2, $3, $4)
			RETURNING preset_id, created_at
		`

		var (
			id      int64
			created time.Time
		)

		if err := database.WithQuery(ctx, r.log, insertPreset, func() error {
			return tx.QueryRowContext(ctx, insertPreset,
				p.Name, p.Description, p.TotalPrice, p.ImageURL,
			).Scan(&id, &created)
		}); err != nil {
			return 0, r.mapPostgreSQLError(err)
		}

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

			if err := database.WithQuery(ctx, r.log, sqlStr, func() error {
				_, execErr := tx.ExecContext(ctx, sqlStr, args...)
				return execErr
			}); err != nil {
				return 0, r.mapPostgreSQLError(err)
			}
		}

		p.ID, p.CreatedAt = id, created
		return id, nil
	})
}

// Get returns preset with embedded Items slice.
func (r *PGPresetRepository) Get(ctx context.Context, id int64) (*preset.Preset, error) {
	const qPreset = `
		SELECT preset_id, name, description, total_price, image_url, created_at
		FROM presets WHERE preset_id = $1
	`

	var raw presetDB
	if err := database.WithQuery(ctx, r.log, qPreset, func() error {
		return r.db.GetContext(ctx, &raw, qPreset, id)
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, preset.ErrPresetNotFound
		}
		return nil, r.mapPostgreSQLError(err)
	}

	p := raw.toDomain()

	const qItems = `
		SELECT preset_item_id, preset_id, product_id
		FROM preset_items WHERE preset_id = $1
	`
	if err := database.WithQuery(ctx, r.log, qItems, func() error {
		return r.db.SelectContext(ctx, &p.Items, qItems, id)
	}); err != nil {
		return nil, r.mapPostgreSQLError(err)
	}

	return p, nil
}

func (r *PGPresetRepository) ListDetailed(ctx context.Context) ([]preset.Preset, error) {
	const qPresets = `
		SELECT preset_id, name, description, total_price, image_url, created_at
		FROM presets
	`

	var raws []presetDB
	if err := database.WithQuery(ctx, r.log, qPresets, func() error {
		return r.db.SelectContext(ctx, &raws, qPresets)
	}); err != nil {
		return nil, r.mapPostgreSQLError(err)
	}

	presets := rawPresetListToDomain(raws)
	if len(presets) == 0 {
		return presets, nil
	}

	ids := make([]int64, len(presets))
	for i := range presets {
		ids[i] = presets[i].ID
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
		WHERE pi.preset_id = ANY($1)
	`

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
		return nil, r.mapPostgreSQLError(err)
	}

	byID := make(map[int64]*preset.Preset, len(presets))
	for i := range presets {
		byID[presets[i].ID] = &presets[i]
	}

	for _, row := range rows {
		item := preset.PresetItem{
			ID:        row.PresetItemID,
			PresetID:  row.PresetID,
			ProductID: row.ProductID,
			Product: &product.ProductSummary{
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

func (r *PGPresetRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM presets WHERE preset_id = $1`
	var res sql.Result
	err := r.withQuery(ctx, q, func() error {
		var execErr error
		res, execErr = r.db.ExecContext(ctx, q, id)
		return execErr
	})
	if err != nil {
		return r.mapPostgreSQLError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return preset.ErrPresetNotFound
	}
	return nil
}

func optionalString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func (r *PGPresetRepository) withQuery(ctx context.Context, query string, fn func() error, extras ...slog.Attr) error {
	return database.WithQuery(ctx, r.log, query, fn, extras...)
}

func (r *PGPresetRepository) mapPostgreSQLError(err error) error {
	return repoError.MapPostgreSQLError(r.log, err)
}