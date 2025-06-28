package attr

import "errors"

var (
	ErrAttributeNotFound   = errors.New("attribute not found")
	ErrAttributeConflict   = errors.New("attribute conflict")
	ErrAttributeValidation = errors.New("attribute validation error")
)
