package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotFound     = errors.New("record not found")
	ErrInvalidRefID = errors.New("referenced record does not exist")
)

func wrapPgError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case pgUniqueViolation:
			return ErrDuplicate
		case pgForeignKeyViolation:
			return ErrInvalidRefID
		}
	}
	return err
}
