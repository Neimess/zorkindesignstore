package dto

//swaggo:model PresetResponse
type PresetResponse struct {
	PresetID    int64                `json:"preset_id" example:"1"`
	Name        string               `json:"name" example:"Комплект для ванной"`
	Description *string              `json:"description,omitempty" example:"Полный комплект для ванной комнаты"`
	TotalPrice  float64              `json:"total_price" example:"15000"`
	ImageURL    *string              `json:"image_url,omitempty" example:"https://example.com/image.png"`
	CreatedAt   string               `json:"created_at" example:"2025-06-20T15:00:00Z"`
	Items       []PresetResponseItem `json:"items"`
}

//swaggo:model PresetShortResponse
type PresetShortResponse struct {
	PresetID    int64   `json:"preset_id" example:"1"`
	Name        string  `json:"name" example:"Комплект для ванной"`
	Description *string `json:"description,omitempty" example:"Для ванной комнаты"`
	TotalPrice  float64 `json:"total_price,omitempty" example:"15000"` // может быть опциональным
	ImageURL    *string `json:"image_url,omitempty" example:"https://example.com/image.png"`
	CreatedAt   string  `json:"created_at" example:"2025-06-20T15:00:00Z"`
}

//swaggo:model PresetResponseItem
type PresetResponseItem struct {
	Product ProductSummary `json:"product"`
}

//swaggo:model ProductSummary
type ProductSummary struct {
	ID       int64   `json:"id" example:"10"`
	Name     string  `json:"name" example:"Шампунь"`
	Price    float64 `json:"price" example:"499"`
	ImageURL *string `json:"image_url,omitempty" example:"https://example.com/shampoo.png"`
}
