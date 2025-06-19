package domain

import "time"

type Product struct {
	ID         int64  `db:"product_id"`
	Name       string `db:"name"`
	Price      int64  `db:"price"`
	CategoryID int64  `db:"category_id"`
}

type ProductImage struct {
	ID        int64  `db:"image_id"`
	ProductID int64  `db:"product_id"`
	URL       string `db:"url"`
	AltText   string `db:"alt_text"`
}

type Preset struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	ImageURL    string    `db:"image_url"`
	ProductIDs  []int64   `db:"-"`
	CreatedAt   time.Time `db:"created_at"`
	IsActive    bool      `db:"is_active"`
}

type PresetProduct struct {
	PresetID  int64 `db:"preset_id"`
	ProductID int64 `db:"product_id"`
}

type ProductWithImages struct {
	Product
	Images []ProductImage
}

type ProductWithCharacteristics struct {
	Product
	Images          []ProductImage
	Characteristics []CategoryCharacteristic
}
