package domain

type Category struct {
	ID   int64  `db:"category_id"`
	Name string `db:"name"`
}

type CategoryAttribute struct {
	CategoryID  int64 `db:"category_id"`
	AttributeID int64 `db:"attribute_id"`
	IsRequired  bool  `db:"is_required"`
	Priority    int   `db:"priority"`
}
