package domain

type Category struct {
	ID   int64  `json:"id" db:"category_id"`
	Name string `json:"name" db:"name"`
}

type CategoryAttribute struct {
    CategoryID  int64 `json:"category_id" db:"category_id"`
    AttributeID int64 `json:"attribute_id" db:"attribute_id"`
    IsRequired  bool  `json:"is_required" db:"is_required"`
}

type CategoryAttributePriority struct {
    CategoryID  int64 `json:"category_id" db:"category_id"`
    AttributeID int64 `json:"attribute_id" db:"attribute_id"`
    Priority    int   `json:"priority" db:"priority"`
}