package product

import (
	"time"
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

type ProductSummary struct {
	ID       int64
	Name     string
	Price    float64
	ImageURL *string
}
