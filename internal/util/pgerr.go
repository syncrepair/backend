package util

import (
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	PgErrNotFound       = errors.New("not found")
	PgErrAlreadyExists  = errors.New("already exists")
	PgErrInternalServer = errors.New("internal server")
)

func ParsePgErr(err error) error {
	var e *pgconn.PgError

	if errors.As(err, &e) {
		if e.Code == pgerrcode.UniqueViolation {
			return PgErrAlreadyExists
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return PgErrNotFound
	}

	return PgErrInternalServer
}
