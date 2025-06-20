package dto
//go:generate easyjson -all
// swagger:model IDResponse
type IDResponse struct {
	ID      int64  `json:"id" example:"42"`
	Message string `json:"message" example:"created"`
}