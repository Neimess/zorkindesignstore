//go:generate easyjson -all
package attribute

// CreateAttributeReq defines the body for creating an attribute
// swagger:model CreateAttributeReq
//
//easyjson:json
type CreateAttributeReq struct {
	Name string  `json:"name" validate:"required"`
	Unit *string `json:"unit,omitempty"`
}

// swagger:model UpdateAttributeReq
//
//easyjson:json
type UpdateAttributeReq struct {
	Name string  `json:"name" validate:"required"`
	Unit *string `json:"unit,omitempty"`
}

// AttributeResponse is the JSON response for attribute endpoints
// swagger:model AttributeResponse
//
//easyjson:json
type AttributeResponse struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Unit       *string `json:"unit,omitempty"`
	CategoryID int64   `json:"categoryID"`
}

// swagger:model AttributeListResponse
//
//easyjson:json
type AttributeListResponse []AttributeResponse

// swagger:model CreateAttributesBatchReq
//
//easyjson:json
type CreateAttributesBatchReq []CreateAttributeReq

// ErrorResponse standard error response
// swagger:model ErrorResponse
//
//easyjson:json
type ErrorResponse struct {
	Error string `json:"error"`
}
