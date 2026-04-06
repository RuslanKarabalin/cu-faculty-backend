package repository

import "errors"

const pgUniqueViolation = "23505"

var (
	ErrDuplicate = errors.New("record already exists")
	ErrNotFound  = errors.New("record not found")
)
