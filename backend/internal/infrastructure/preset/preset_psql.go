package preset

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/preset"
	repoError "github.com/Neimess/zorkin-store-project/internal/infrastructure/error"
	"github.com/Neimess/zorkin-store-project/pkg/database"
	"github.com/Neimess/zorkin-store-project/pkg/database/tx"
	logger "github.com/Neimess/zorkin-store-project/pkg/log"
)

const (
	opCreate       = "repo.preset.Create"
	opUpdate       = "repo.preset.Update"
	createdAtField = "created_at"
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

func (r *PGPresetRepository) Create(ctx context.Context, p *preset.Preset) (*preset.Preset, error) {
	r.log.Debug("Create preset", slog.String("op", opCreate), slog.String("name", p.Name))
	return r.save(ctx, p, true)
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
		return nil, r.mapPostgreSQLError(err)
	}

	p := raw.toDomain()

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
		WHERE pi.preset_id = $1
	`

	var rows []presetItemDetailedDB

	if err := database.WithQuery(ctx, r.log, qItems, func() error {
		return r.db.SelectContext(ctx, &rows, qItems, id)
	}); err != nil {
		return nil, r.mapPostgreSQLError(err)
	}
	for _, row := range rows {
		item := row.toDomain()
		p.Items = append(p.Items, item)
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

	var rows []presetItemDetailedDB
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
		item := row.toDomain()
		byID[row.PresetID].Items = append(byID[row.PresetID].Items, item)
	}

	return presets, nil
}

func (r *PGPresetRepository) ListShort(ctx context.Context) ([]preset.Preset, error) {
	const q = `
		SELECT preset_id, name, description, total_price, image_url, created_at
		FROM presets
	`
	var raws []presetDB
	if err := r.withQuery(ctx, q, func() error {
		return r.db.SelectContext(ctx, &raws, q)
	}); err != nil {
		return nil, r.mapPostgreSQLError(err)
	}

	return rawPresetListToDomain(raws), nil
}

func (r *PGPresetRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM presets WHERE preset_id=$1`
	err := database.WithQuery(ctx, r.log, q, func() error {
		_, execErr := r.db.ExecContext(ctx, q, id)
		return execErr
	})
	if err != nil {
		return r.mapPostgreSQLError(err)
	}
	return nil
}

func (r *PGPresetRepository) Update(ctx context.Context, p *preset.Preset) (*preset.Preset, error) {
	r.log.Debug("Update preset", slog.String("op", opUpdate), slog.Int64("id", p.ID))
	return r.save(ctx, p, false)
}

func (r *PGPresetRepository) save(ctx context.Context, p *preset.Preset, isNew bool) (*preset.Preset, error) {
	queryPreset := `UPDATE presets SET name=$1, description=$2, total_price=$3, image_url=$4 WHERE preset_id=$5`
	if isNew {
		queryPreset = `INSERT INTO presets (name, description, total_price, image_url) VALUES ($1,$2,$3,$4) RETURNING preset_id, created_at`
	}

	resPreset, err := tx.RunInTx(ctx, r.db, func(tx *sqlx.Tx) (*preset.Preset, error) {
		// Сохранение Preset
		err := database.WithQuery(ctx, r.log, queryPreset, func() error {
			if isNew {
				return tx.QueryRowContext(ctx, queryPreset, p.Name, p.Description, p.TotalPrice, p.ImageURL).
					Scan(&p.ID, &p.CreatedAt)
			}
			_, execErr := tx.ExecContext(ctx, queryPreset, p.Name, p.Description, p.TotalPrice, p.ImageURL, p.ID)
			return execErr
		})
		if err != nil {
			return nil, r.mapPostgreSQLError(err)
		}

		// Перезапись элементов
		err = r.deleteItems(ctx, tx, p.ID)
		if err != nil {
			return nil, err
		}

		if len(p.Items) > 0 {
			err = r.insertItems(ctx, tx, p.ID, p.Items)
			if err != nil {
				return nil, err
			}
		}
		return p, nil
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, preset.ErrPresetNotFound
		}
		return nil, err
	}
	return resPreset, nil
}

func (r *PGPresetRepository) deleteItems(ctx context.Context, tx *sqlx.Tx, presetID int64) error {
	const q = `DELETE FROM preset_items WHERE preset_id=$1`
	err := database.WithQuery(ctx, r.log, q, func() error {
		_, execErr := tx.ExecContext(ctx, q, presetID)
		return execErr
	})
	if err != nil {
		return r.mapPostgreSQLError(err)
	}
	return nil
}

func (r *PGPresetRepository) insertItems(ctx context.Context, tx *sqlx.Tx, presetID int64, items []preset.PresetItem) error {
	n := len(items)
	ids := make([]int64, n)
	pids := make([]int64, n)
	for i, it := range items {
		ids[i] = presetID
		pids[i] = it.ProductID
	}
	const q = `INSERT INTO preset_items (preset_id, product_id) SELECT * FROM UNNEST($1::bigint[],$2::bigint[])`
	err := database.WithQuery(ctx, r.log, q, func() error {
		_, execErr := tx.ExecContext(ctx, q, pq.Array(ids), pq.Array(pids))
		return execErr
	})
	if err != nil {
		return r.mapPostgreSQLError(err)
	}
	return nil
}

func (r *PGPresetRepository) withQuery(ctx context.Context, query string, fn func() error, extras ...slog.Attr) error {
	return database.WithQuery(ctx, r.log, query, fn, extras...)
}

func (r *PGPresetRepository) mapPostgreSQLError(err error) error {
	return repoError.MapPostgreSQLError(r.log, err)
}
