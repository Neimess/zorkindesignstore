package attribute

import (
	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
)

type rawSQL struct {
	ID         int64   `db:"attribute_id"`
	Name       string  `db:"name"`
	Unit       *string `db:"unit"`
	CategoryID int64   `db:"category_id"`
}

func (r rawSQL) toDomain() *attr.Attribute {
	return &attr.Attribute{
		ID:         r.ID,
		Name:       r.Name,
		Unit:       r.Unit,
		CategoryID: r.CategoryID,
	}
}

func rawListToDomain(raws []rawSQL) []attr.Attribute {
	attrs := make([]attr.Attribute, len(raws))
	for i, r := range raws {
		attrs[i] = *r.toDomain()
	}
	return attrs
}
