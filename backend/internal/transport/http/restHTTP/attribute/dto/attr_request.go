package dto

import (
	attr "github.com/Neimess/zorkin-store-project/internal/domain/attribute"
)

type AttributeRequest struct {
	Name string  `json:"name" validate:"required,min=1,max=255"`
	Unit *string `json:"unit,omitempty" validate:"omitempty,max=50"`
}

type CreateAttributesBatchRequest struct {
	Items []AttributeBatchDataItem `json:"data" validate:"required,min=1,dive"`
}

type AttributeBatchDataItem struct {
	Name string  `json:"name" validate:"required,min=1,max=255"`
	Unit *string `json:"unit,omitempty" validate:"omitempty,max=50"`
}

func (r *AttributeRequest) MapToDomain() *attr.Attribute {
	return &attr.Attribute{
		Name: r.Name,
		Unit: r.Unit,
	}
}

func (r *CreateAttributesBatchRequest) MapToDomainBatch() []attr.Attribute {
	out := make([]attr.Attribute, len(r.Items))
	for i, it := range r.Items {
		out[i] = attr.Attribute{
			Name: it.Name,
			Unit: it.Unit,
		}
	}
	return out
}
