package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const pgUniqueViolation = "23505"

var (
	ErrDuplicate = errors.New("record already exists")
	ErrNotFound  = errors.New("record not found")
)

func wrapPgError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == pgUniqueViolation {
		return ErrDuplicate
	}
	return err
}
