package repository

import (
	"errors"
)

const (
	foreignKeyViolationCode string = "23503"
	uniqueViolationCode     string = "23505"
)

var (
	ErrInvalidForeignKey = errors.New("invalid attribute_id: does not exist")
	ErrDuplicateName	 = errors.New("duplicate name: already exists")

)
