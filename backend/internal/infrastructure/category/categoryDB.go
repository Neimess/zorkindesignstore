package category

type categoryDB struct {
	ID   int64  `db:"category_id"`
	Name string `db:"name"`
}
