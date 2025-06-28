package dto

import "github.com/Neimess/zorkin-store-project/internal/domain/preset"

//swaggo:model PresetRequest
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

func (r *PresetRequest) MapToPreset() *preset.Preset {
	return &preset.Preset{
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