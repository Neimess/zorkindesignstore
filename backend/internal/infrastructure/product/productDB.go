package product

import (
	"database/sql"
	"time"

	attrDom "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
	prodDom "github.com/Neimess/zorkin-store-project/internal/domain/product"
)

type productRow struct {
	ID          int64          `db:"product_id"`
	Name        string         `db:"name"`
	Price       float64        `db:"price"`
	Description sql.NullString `db:"description"`
	CategoryID  int64          `db:"category_id"`
	ImageURL    sql.NullString `db:"image_url"`
	CreatedAt   time.Time      `db:"created_at"`
}

func (r *productRow) toDomain(attrs []prodDom.ProductAttribute) *prodDom.Product {
	d := &prodDom.Product{
		ID:          r.ID,
		Name:        r.Name,
		Price:       r.Price,
		Description: r.Description.String,
		CategoryID:  r.CategoryID,
		ImageURL:    r.ImageURL.String,
		CreatedAt:   r.CreatedAt,
		Attributes:  attrs,
	}
	return d
}

type productAttributeRow struct {
	ProductID   int64          `db:"product_id"`
	AttributeID int64          `db:"attribute_id"`
	Value       string         `db:"value"`
	Name        string         `db:"name"`
	Unit        sql.NullString `db:"unit"`
}

func (ar *productAttributeRow) toDomainAttr() prodDom.ProductAttribute {
	var unit *string
	if ar.Unit.Valid {
		unit = &ar.Unit.String
	}
	return prodDom.ProductAttribute{
		ProductID:   ar.ProductID,
		AttributeID: ar.AttributeID,
		Value:       ar.Value,
		Attribute: attrDom.Attribute{
			ID:   ar.AttributeID,
			Name: ar.Name,
			Unit: unit,
		},
	}
}
