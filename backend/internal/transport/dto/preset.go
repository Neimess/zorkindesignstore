package dto

import (
	"time"

	"github.com/Neimess/zorkin-store-project/internal/domain"
)

// -----------------------------------------------------------------------------
// ❖ REQUESTS
// -----------------------------------------------------------------------------

//swaggo:model PresetRequest
//go:generate easyjson -all
type PresetRequest struct {
	Name        string              `json:"name" validate:"required,min=2"`
	Description string              `json:"description,omitempty"`
	TotalPrice  float64             `json:"total_price" validate:"gt=0"`
	ImageURL    string              `json:"image_url,omitempty" validate:"omitempty,url"`
	Items       []PresetRequestItem `json:"items" validate:"dive,required"`
}

//swaggo:model PresetRequestItem
type PresetRequestItem struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
}

// -----------------------------------------------------------------------------
// ❖ RESPONSES
// -----------------------------------------------------------------------------

// ProductSummary — облегчённая карточка продукта, которую мы вкладываем в PresetItem.
//
//swaggo:model ProductSummary
//easyjson:json
type ProductSummary struct {
	ID       int64   `json:"id" example:"10"`
	Name     string  `json:"name" example:"Шампунь"`
	Price    float64 `json:"price" example:"499"`
	ImageURL *string `json:"image_url,omitempty" example:"https://example.com/shampoo.png"`
}

//swaggo:model PresetResponse
//easyjson:json
type PresetResponse struct {
	PresetID    int64                `json:"preset_id" example:"1"`
	Name        string               `json:"name" example:"Комплект для ванной"`
	Description string               `json:"description,omitempty" example:"Полный комплект для ванной комнаты"`
	TotalPrice  float64              `json:"total_price" example:"15000"`
	ImageURL    string               `json:"image_url,omitempty" example:"https://example.com/image.png"`
	CreatedAt   string               `json:"created_at" example:"2025-06-20T15:00:00Z"`
	Items       []PresetResponseItem `json:"items"`
}

//swaggo:model PresetShortResponse
//easyjson:json
type PresetShortResponse struct {
	PresetID    int64   `json:"preset_id" example:"1"`
	Name        string  `json:"name" example:"Комплект для ванной"`
	Description string  `json:"description,omitempty" example:"Для ванной комнаты"`
	TotalPrice  float64 `json:"total_price,omitempty" example:"15000"` // может быть опциональным
	ImageURL    string  `json:"image_url,omitempty" example:"https://example.com/image.png"`
	CreatedAt   string  `json:"created_at" example:"2025-06-20T15:00:00Z"`
}

//swaggo:model PresetResponseItem
//easyjson:json
type PresetResponseItem struct {
	Product ProductSummary `json:"product"`
}

// easyjson:json
type PresetListResponse []PresetResponse

// easyjson:json
type PresetShortListResponse []PresetShortResponse

// -----------------------------------------------------------------------------
// ❖ MAPPERS
// -----------------------------------------------------------------------------

// MapToPreset конвертирует DTO‑запрос в доменную модель Preset.
func (r *PresetRequest) MapToPreset() *domain.Preset {
	return &domain.Preset{
		Name:        r.Name,
		Description: r.Description,
		TotalPrice:  r.TotalPrice,
		ImageURL:    r.ImageURL,
		Items:       r.MapToPresetItems(),
	}
}

func MapToShortDto(p *domain.Preset) PresetShortResponse {
	return PresetShortResponse{
		PresetID:    p.ID,
		Name:        p.Name,
		Description: p.Description,
		TotalPrice:  p.TotalPrice,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
	}
}

func (r *PresetRequest) MapToPresetItems() []domain.PresetItem {
	items := make([]domain.PresetItem, len(r.Items))
	for i, it := range r.Items {
		items[i] = domain.PresetItem{ProductID: it.ProductID}
	}
	return items
}

// MapToDto делает из доменной Preset DTO для ответа.
func MapToDto(p *domain.Preset) *PresetResponse {
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

func mapToResponseItems(items []domain.PresetItem) []PresetResponseItem {
	out := make([]PresetResponseItem, len(items))
	for i, it := range items {
		ps := ProductSummary{ID: it.ProductID}
		if it.Product != nil { // если репозиторий заполнил краткую инфу
			ps.Name = it.Product.Name
			ps.Price = it.Product.Price
			ps.ImageURL = it.Product.ImageURL
		}
		out[i] = PresetResponseItem{Product: ps}
	}
	return out
}
