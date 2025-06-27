package attribute

type CreateAttributeInput struct {
	Name       string  `validate:"required,min=1,max=100"`
	Unit       *string `validate:"omitempty,max=50"`
	CategoryID int64   `validate:"required,gte=1"`
}

type CreateAttributesBatchInput []*CreateAttributeInput
