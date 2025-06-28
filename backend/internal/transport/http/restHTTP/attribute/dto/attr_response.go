package dto

import (
	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
)

type AttributeResponse struct {
	ID         int64   `json:"id"` // attribute_id
	Name       string  `json:"name"`
	Unit       *string `json:"unit,omitempty"`
	CategoryID int64   `json:"category_id"` // к какой категории относится
}

func MapToAttributeResponse(a *attr.Attribute) AttributeResponse {
	return AttributeResponse{
		ID:         a.ID,
		Name:       a.Name,
		Unit:       a.Unit,
		CategoryID: a.CategoryID,
	}
}

type AttributeListResponse []AttributeResponse

func MapToAttributeListResponse(attrs []*attr.Attribute) AttributeListResponse {
	out := make(AttributeListResponse, len(attrs))
	for i, a := range attrs {
		out[i] = MapToAttributeResponse(a)
	}
	return out
}
