package domain

import "errors"

var (
	ErrDuplicate   = errors.New("already exists")
	ErrInvalidData = errors.New("invalid data")
	ErrDatabase    = errors.New("database error")
)
