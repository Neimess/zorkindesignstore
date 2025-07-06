package product

import (
	"time"

	serviceDom "github.com/Neimess/zorkin-store-project/internal/domain/service"
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
	Services    []serviceDom.Service
}

type ProductSummary struct {
	ID       int64
	Name     string
	Price    float64
	ImageURL *string
}
