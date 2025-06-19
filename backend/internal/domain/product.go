package domain

import "time"


type Product struct {
    ID          int64     `json:"id" db:"product_id"`
    Name        string    `json:"name" db:"name"`
    Price       float64   `json:"price" db:"price"`
    Description string    `json:"description,omitempty" db:"description"`
    CategoryID  int64     `json:"category_id" db:"category_id"`
    ImageURL    string    `json:"image_url,omitempty" db:"image_url"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type ProductAttribute struct {
    ProductID   int64    `json:"product_id" db:"product_id"`
    AttributeID int64    `json:"attribute_id" db:"attribute_id"`

    ValueString *string  `json:"value_string,omitempty" db:"value_string"`
    ValueInt    *int64   `json:"value_int,omitempty" db:"value_int"`
    ValueFloat  *float64 `json:"value_float,omitempty" db:"value_float"`
    ValueBool   *bool    `json:"value_bool,omitempty" db:"value_bool"`
    ValueEnum   *string  `json:"value_enum,omitempty" db:"value_enum"`
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
