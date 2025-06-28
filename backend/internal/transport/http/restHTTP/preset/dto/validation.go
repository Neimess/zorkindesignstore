package dto

type ValidationErrorResponse struct {
    Message string       `json:"message"`
    Errors  []FieldError `json:"errors"`
}

type FieldError struct {
    Field string `json:"field"`
    Tag   string `json:"tag"`
    Value string `json:"value"`
}
