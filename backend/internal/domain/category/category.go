package category

import (
	"errors"
	"strings"
)

const maxNameLength = 255

var (
	ErrCategoryNameEmpty   = errors.New("category name cannot be empty")
	ErrCategoryNameTooLong = errors.New("category name must be at most 255 characters")
	ErrAttributeNameEmpty  = errors.New("attribute name cannot be empty")
	ErrCategoryNotFound    = errors.New("category not found")
	ErrCategoryInUse       = errors.New("category is in use and cannot be deleted")
)

type Category struct {
	ID   int64
	Name string
}

// Validate проверяет бизнес-правила сущности Category:
// 1. Обрезает пробелы вокруг Name.
// 2. Гарантирует, что Name не пустой.
// 3. Проверяет, что длина Name не превышает maxNameLength.
func (c *Category) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return ErrCategoryNameEmpty
	}
	if len(c.Name) > maxNameLength {
		return ErrCategoryNameTooLong
	}
	return nil
}
