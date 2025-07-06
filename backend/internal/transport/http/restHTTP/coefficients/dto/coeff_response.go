package dto

type CoefficientResponse struct {
	ID    int64   `json:"id" example:"1"`
	Name  string  `json:"name" example:"Коэффициент 1"`
	Value float64 `json:"value" example:"1.2345"`
}
