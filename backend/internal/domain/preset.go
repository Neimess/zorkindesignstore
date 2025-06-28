package domain

import (
	"time"

	"github.com/Neimess/zorkin-store-project/internal/domain/product"
)

type Preset struct {
	ID          int64        `json:"preset_id" db:"preset_id"`
	Name        string       `json:"name" db:"name"`
	Description string       `json:"description,omitempty" db:"description"`
	TotalPrice  float64      `json:"total_price" db:"total_price"`
	ImageURL    string       `json:"image_url,omitempty" db:"image_url"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	Items       []PresetItem `json:"items,omitempty" db:"-"`
}

type PresetItem struct {
	ID        int64                   `json:"id" db:"preset_item_id"`
	PresetID  int64                   `json:"preset_id" db:"preset_id"`
	ProductID int64                   `json:"product_id" db:"product_id"`
	Product  *product.ProductSummary `json:"product,omitempty"`
}
