package domain

import "errors"

// Errors used in the domain layer of the application.
var (
	ErrNotFound = errors.New("entity not found")
)
