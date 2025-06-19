package domain

type AttrDataType string

const (
	AttrDataTypeString AttrDataType = "string"
	AttrDataTypeInt    AttrDataType = "int"
	AttrDataTypeFloat  AttrDataType = "float"
	AttrDataTypeBool   AttrDataType = "bool"
	AttrDataTypeEnum   AttrDataType = "enum"
)

//go:generate easyjson -all
type Attribute struct {
	ID           int64        `json:"id" db:"attribute_id"`
	Name         string       `json:"name" db:"name"`
	Slug         string       `json:"slug" db:"slug"`
	DataType     AttrDataType `json:"data_type" db:"data_type"`
	Unit         string      `json:"unit,omitempty" db:"unit"`
	IsFilterable bool         `json:"is_filterable" db:"is_filterable"`
}
