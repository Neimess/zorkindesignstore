package preset

import (
	"database/sql"
	"time"

	presetDom "github.com/Neimess/zorkin-store-project/internal/domain/preset"
	"github.com/Neimess/zorkin-store-project/internal/domain/product"
)

type presetDB struct {
	ID          int64          `db:"preset_id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	TotalPrice  float64        `db:"total_price"`
	ImageURL    sql.NullString `db:"image_url"`
	CreatedAt   time.Time      `db:"created_at"`
}

func (p presetDB) toDomain() *presetDom.Preset {
	return &presetDom.Preset{
		ID:          p.ID,
		Name:        p.Name,
		Description: optionalString(p.Description),
		TotalPrice:  p.TotalPrice,
		ImageURL:    optionalString(p.ImageURL),
		CreatedAt:   p.CreatedAt,
	}
}

// rawPresetListToDomain массово преобразует список.
func rawPresetListToDomain(raws []presetDB) []presetDom.Preset {
	presets := make([]presetDom.Preset, len(raws))
	for i, r := range raws {
		presets[i] = *r.toDomain()
	}
	return presets
}

type presetItemDetailedDB struct {
	ID           int64          `db:"preset_item_id"`
	PresetID     int64          `db:"preset_id"`
	ProductID    int64          `db:"product_id"`
	ProductName  string         `db:"product_name"`
	ProductPrice float64        `db:"product_price"`
	ProductImage sql.NullString `db:"product_image_url"`
}

func (r presetItemDetailedDB) toDomain() presetDom.PresetItem {
	return presetDom.PresetItem{
		ID:        r.ID,
		PresetID:  r.PresetID,
		ProductID: r.ProductID,
		Product: &product.ProductSummary{
			ID:       r.ProductID,
			Name:     r.ProductName,
			Price:    r.ProductPrice,
			ImageURL: optionalString(r.ProductImage),
		},
	}
}

func optionalString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}
