package product

import "errors"

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrBadCategoryID    = errors.New("invalid category ID")
	ErrInvalidAttribute = errors.New("invalid product attribute data")
)
