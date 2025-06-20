package dto

// swagger:model ErrorResponse
//
//go:generate easyjson -all
type ErrorResponse struct {
	Message string `json:"message" example:"invalid request payload"`
}
