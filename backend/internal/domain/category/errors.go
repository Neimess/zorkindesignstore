package category

import "errors"

var (
	ErrCategoryNameEmpty   = errors.New("category name cannot be empty")
	ErrCategoryNameTooLong = errors.New("category name must be at most 255 characters")
	ErrAttributeNameEmpty  = errors.New("attribute name cannot be empty")
	ErrCategoryNotFound    = errors.New("category not found")
	ErrCategoryInUse       = errors.New("category is in use and cannot be deleted")
	ErrCategoryNameExists  = errors.New("category name already exists")
)
