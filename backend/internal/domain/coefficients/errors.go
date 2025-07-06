package coefficients

import "errors"

var (
	ErrCoefficientNotFound      = errors.New("coefficient not found")
	ErrCoefficientAlreadyExists = errors.New("coefficient already exists")
	ErrEmptyName                = errors.New("coefficient name must not be empty")
	ErrNameTooLong              = errors.New("coefficient name is too long")
)
