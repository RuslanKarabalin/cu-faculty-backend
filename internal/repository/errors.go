package repository

import "errors"

const pgUniqueViolation = "23505"

var (
	ErrDuplicate   = errors.New("record already exists")
	ErrNotFound    = errors.New("record not found")
	ErrInvalidDate = errors.New("invalid birth_date format, expected YYYY-MM-DD")
)
