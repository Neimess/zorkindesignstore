package coefficients

import (
	"strings"
)

type Coefficient struct {
	ID    int64
	Name  string
	Value float64
}

func (c *Coefficient) Validate() error {
	name := strings.TrimSpace(c.Name)
	if name == "" {
		return ErrEmptyName
	}
	if len(name) > 255 {
		return ErrNameTooLong
	}
	return nil
}
