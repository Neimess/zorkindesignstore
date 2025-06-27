package domain

import (
	"errors"
)

var (
	ErrCategoryNameEmpty  = errors.New("category name cannot be empty")
	ErrAttributeNameEmpty = errors.New("attribute name cannot be empty")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCategoryInUse      = errors.New("category is in use and cannot be deleted")
)

type Category struct {
	ID   int64
	Name string
}
