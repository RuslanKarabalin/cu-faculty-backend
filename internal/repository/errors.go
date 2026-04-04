package repository

import "errors"

var (
	ErrDuplicate   = errors.New("record already exists")
	ErrInvalidDate = errors.New("invalid birth_date format, expected YYYY-MM-DD")
)
