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
