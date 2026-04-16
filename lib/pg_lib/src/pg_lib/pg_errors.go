package pg_lib

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func IsNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func IsConstraintViolated(err *pgconn.PgError) bool {
	return err.ConstraintName != ""
}

func IsSerializationError(err *pgconn.PgError) bool {
	return err.Code == "40001"
}

func IsServerSideTimeout(err *pgconn.PgError) bool {
	return err.Code == "57014"
}
