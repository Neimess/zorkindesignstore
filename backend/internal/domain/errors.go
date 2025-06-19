package domain

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

func ErrValidation(msg string) error {
	return ValidationError{Message: msg}
}

// IsValidationError checks whether the error is a ValidationError.
func IsValidationError(err error) bool {
	var ve ValidationError
	return errors.As(err, &ve)
}
