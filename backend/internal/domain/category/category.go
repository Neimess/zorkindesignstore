package category

import (
	"strings"
)

const maxNameLength = 255

type Category struct {
	ID   int64
	Name string
}

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
