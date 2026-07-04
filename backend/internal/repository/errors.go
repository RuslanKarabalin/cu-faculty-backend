package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgUniqueViolation           = "23505"
	pgForeignKeyViolation       = "23503"
	pgNotNullViolation          = "23502"
	pgCheckViolation            = "23514"
	pgStringDataRightTruncation = "22001"
	pgNumericValueOutOfRange    = "22003"
	pgInvalidTextRepresentation = "22P02"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotFound     = errors.New("record not found")
	ErrInvalidRefID = errors.New("referenced record does not exist")
	ErrValidation   = errors.New("validation failed")
)

func wrapPgError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case pgUniqueViolation:
			return ErrDuplicate
		case pgForeignKeyViolation:
			return ErrInvalidRefID
		case pgNotNullViolation,
			pgCheckViolation,
			pgStringDataRightTruncation,
			pgNumericValueOutOfRange,
			pgInvalidTextRepresentation:
			return ErrValidation
		}
	}
	return err
}
