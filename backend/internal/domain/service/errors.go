package service

import "errors"

var (
	ErrServiceNotFound      = errors.New("service not found")
	ErrServiceAlreadyExists = errors.New("service already exists")
	ErrEmptyName            = errors.New("service name must not be empty")
	ErrNameTooLong          = errors.New("service name is too long")
)
