package dto

// swagger:model IDResponse
//
//go:generate easyjson -all
type IDResponse struct {
	ID      int64  `json:"id" example:"42"`
	Message string `json:"message" example:"created"`
}
