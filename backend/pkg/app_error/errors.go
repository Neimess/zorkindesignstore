package app_error

import "errors"

var (
	ErrConflict   = errors.New("conflict")
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation error")
	ErrInternal   = errors.New("internal server error")
	ErrCanceled   = errors.New("operation canceled")
	ErrTimeout    = errors.New("operation timeout")
)

// 23505
