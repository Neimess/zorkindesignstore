package dto

type ServiceResponse struct {
	ID          int64   `json:"id" example:"1"`
	Name        string  `json:"name" example:"Монтаж"`
	Description *string `json:"description,omitempty" example:"Установка изделия"`
	Price       float64 `json:"price" example:"1500.00"`
}
