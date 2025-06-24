package domain

type Category struct {
	ID         int64
	Name       string
	Attributes []CategoryAttribute
}

type CategoryAttribute struct {
	IsRequired bool      `db:"is_required"`
	Priority   int       `db:"priority"`
	Attr       Attribute `db:"-"`
}

func NewCategory(name string, attrs []CategoryAttribute) *Category {
	return &Category{
		Name:       name,
		Attributes: attrs,
	}
}
