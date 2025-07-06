package category

import "database/sql"

type categoryDB struct {
	ID       int64         `db:"category_id"`
	Name     string        `db:"name"`
	ParentID sql.NullInt64 `db:"parent_id"`
}
