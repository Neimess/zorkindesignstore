package dto

//go:generate easyjson -all
type TokenResponse struct {
	Token string `json:"token"`
}
