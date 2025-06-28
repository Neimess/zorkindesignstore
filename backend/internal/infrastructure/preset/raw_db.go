package preset

import (
	"time"

	domain "github.com/Neimess/zorkin-store-project/internal/domain/preset"
)

// presetDB — SQL-структура для таблицы `presets`.
type presetDB struct {
	ID          int64     `db:"preset_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	TotalPrice  float64   `db:"total_price"`
	ImageURL    *string   `db:"image_url"`
	CreatedAt   time.Time `db:"created_at"`
}

func (p presetDB) toDomain() *domain.Preset {
	return &domain.Preset{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		TotalPrice:  p.TotalPrice,
		ImageURL:    *valueOrNil(p.ImageURL),
		CreatedAt:   p.CreatedAt,
	}
}

type presetItemDB struct {
	ID        int64 `db:"preset_item_id"`
	PresetID  int64 `db:"preset_id"`
	ProductID int64 `db:"product_id"`
}

// rawPresetListToDomain массово преобразует список.
func rawPresetListToDomain(raws []presetDB) []domain.Preset {
	presets := make([]domain.Preset, len(raws))
	for i, r := range raws {
		presets[i] = *r.toDomain()
	}
	return presets
}

func valueOrNil(v *string) *string {
	if v == nil || *v == "" {
		return nil
	}
	return v
}
