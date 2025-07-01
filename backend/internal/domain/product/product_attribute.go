package product

import (
	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
)

type ProductAttribute struct {
	ProductID   int64
	AttributeID int64
	Value       string
	Attribute   attr.Attribute
}
