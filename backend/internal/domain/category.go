package domain

type Category struct {
	ID   int64  `db:"category_id"`
	Name string `db:"category_name"`
}

type CategoryCharacteristic struct {
	ID         int64  `db:"characteristic_id"`
	Name       string `db:"characteristic_name"`
	CategoryID int64  `db:"category_id"`
}
