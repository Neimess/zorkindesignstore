package product

import (
	"time"

	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
)

type Product struct {
	ID          int64
	Name        string
	Price       float64
	Description *string
	CategoryID  int64
	ImageURL    *string
	CreatedAt   time.Time
	Attributes  []ProductAttribute
}

type ProductAttribute struct {
	ProductID   int64
	AttributeID int64
	Value       string
	Attribute   attr.Attribute
}

type PresetProduct struct {
	PresetID  int64
	ProductID int64
}

type ProductSummary struct {
	ID       int64
	Name     string
	Price    float64
	ImageURL *string
}
