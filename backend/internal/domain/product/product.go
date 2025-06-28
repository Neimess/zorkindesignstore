package product

import (
	"time"

	"github.com/Neimess/zorkin-store-project/internal/domain/attribute"
)

type Product struct {
	ID          int64              `db:"product_id"`
	Name        string             `db:"name"`
	Price       float64            `db:"price"`
	Description *string            `db:"description"`
	CategoryID  int64              `db:"category_id"`
	ImageURL    *string            `db:"image_url"`
	CreatedAt   time.Time          `db:"created_at"`
	Attributes  []ProductAttribute `db:"-"`
}

type ProductAttribute struct {
	ProductID   int64          `db:"product_id"`
	AttributeID int64          `db:"attribute_id"`
	Value       string         `db:"value"`
	Attribute   attr.Attribute `db:"-"`
}

type PresetProduct struct {
	PresetID  int64 `db:"preset_id"`
	ProductID int64 `db:"product_id"`
}

type ProductSummary struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	ImageURL *string `json:"image_url,omitempty"`
}
