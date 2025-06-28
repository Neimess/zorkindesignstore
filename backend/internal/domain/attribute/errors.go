package attr

import "errors"

var (
	ErrAttributeNotFound      = errors.New("attribute not found")
	ErrAttributeConflict      = errors.New("attribute conflict")
	ErrAttributeValidation    = errors.New("attribute validation error")
	ErrAttributeAlreadyExists = errors.New("attribute already exists")
	ErrBatchEmpty             = errors.New("batch of attributes is empty")
	ErrBatchTooLarge          = errors.New("batch of attributes is too large")
)
