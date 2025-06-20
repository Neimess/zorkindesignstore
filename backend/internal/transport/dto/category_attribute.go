package dto

//go:generate easyjson -all category_attribute.go
type CategoryAttributeRequest struct {
	CategoryID  int64 `json:"category_id"`
	AttributeID int64 `json:"attribute_id"`
	IsRequired  bool  `json:"is_required"`
	Priority    int   `json:"priority"`
}

//easyjson:json
type CategoryAttributeListResponse []CategoryAttributeResponse
//easyjson:json
type CategoryAttributeResponse struct {
	CategoryID  int64 `json:"category_id"`
	AttributeID int64 `json:"attribute_id"`
	IsRequired  bool  `json:"is_required"`
	Priority    int   `json:"priority"`
}