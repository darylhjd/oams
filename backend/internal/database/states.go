package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type SQLState string

const (
	SQLStateForeignKeyViolation    SQLState = "23503"
	SQLStateDuplicateKeyOrIndex    SQLState = "23505"
	SQLStateCheckConstraintFailure SQLState = "23514"
)

// ErrSQLState is a helper function to check if a given err corresponds to a *pgconn.PgError
// and that the error is of a given SQLState.
func ErrSQLState(err error, wantState SQLState) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.SQLState() == string(wantState)
}
