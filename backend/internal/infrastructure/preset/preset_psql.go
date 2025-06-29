package preset

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/domain/preset"
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
			n := len(p.Items)
			presetIDs := make([]int64, n)
			productIDs := make([]int64, n)

			for i, it := range p.Items {
				presetIDs[i] = id
				productIDs[i] = it.ProductID
			}

			const insertItems = `
				INSERT INTO preset_items (preset_id, product_id)
				SELECT * FROM UNNEST($1::bigint[], $2::bigint[])
			`

			if err := database.WithQuery(ctx, r.log, insertItems, func() error {
				_, execErr := tx.ExecContext(ctx, insertItems, pq.Array(presetIDs), pq.Array(productIDs))
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

func (r *PGPresetRepository) withQuery(ctx context.Context, query string, fn func() error, extras ...slog.Attr) error {
	return database.WithQuery(ctx, r.log, query, fn, extras...)
}

func (r *PGPresetRepository) mapPostgreSQLError(err error) error {
	return repoError.MapPostgreSQLError(r.log, err)
}
