package dto

import (
	"time"

	"github.com/Neimess/zorkin-store-project/internal/domain/preset"
)

// MapToPreset конвертирует DTO‑запрос в доменную модель Preset.

func MapDomainToShortDTO(p *preset.Preset) PresetShortResponse {
	return PresetShortResponse{
		PresetID:    p.ID,
		Name:        p.Name,
		Description: p.Description,
		TotalPrice:  p.TotalPrice,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
	}
}

func MapDomainToDto(p *preset.Preset) *PresetResponse {
	return &PresetResponse{
		PresetID:    p.ID,
		Name:        p.Name,
		Description: p.Description,
		TotalPrice:  p.TotalPrice,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		Items:       mapToResponseItems(p.Items),
	}
}

func mapToResponseItems(items []preset.PresetItem) []PresetResponseItem {
	out := make([]PresetResponseItem, len(items))
	for i, it := range items {
		ps := ProductSummary{ID: it.ProductID}
		if it.Product != nil {
			ps.Name = it.Product.Name
			ps.Price = it.Product.Price
			ps.ImageURL = it.Product.ImageURL
		}
		out[i] = PresetResponseItem{Product: ps}
	}
	return out
}

func (r *PresetRequest) MapToPreset() *preset.Preset {
	return &preset.Preset{
		Name:        r.Name,
		Description: r.Description,
		TotalPrice:  r.TotalPrice,
		ImageURL:    r.ImageURL,
		Items:       r.mapToPresetItems(),
	}
}

func (r *PresetRequest) MapUpdateToPreset(preset_id int64) *preset.Preset {
	return &preset.Preset{
		ID:          preset_id,
		Name:        r.Name,
		Description: r.Description,
		TotalPrice:  r.TotalPrice,
		ImageURL:    r.ImageURL,
		Items:       r.mapToPresetItems(),
	}
}
func (r *PresetRequest) mapToPresetItems() []preset.PresetItem {
	items := make([]preset.PresetItem, len(r.Items))
	for i, it := range r.Items {
		items[i] = preset.PresetItem{ProductID: it.ProductID}
	}
	return items
}
