package domain

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrDuplicate   = errors.New("already exists")
	ErrInvalidData = errors.New("invalid data")
	ErrDatabase    = errors.New("database error")
)
