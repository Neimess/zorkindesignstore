package domain

type Attribute struct {
	ID           int64   `db:"attribute_id"`
	Name         string  `db:"name"`
	Slug         string  `db:"slug"`
	Unit         *string `db:"unit"`
	IsFilterable bool    `db:"is_filterable"`
}
