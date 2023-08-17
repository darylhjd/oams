package database

import . "github.com/go-jet/jet/v2/postgres"

const (
	ListDefaultLimit = 50
)

type limitOffsetter interface {
	limiter
	offsetter
}

type limiter interface {
	Limit() *int64
}

// setLimit sets a limit to the number of returned rows in a select statement.
func setLimit(stmt SelectStatement, limit limiter) SelectStatement {
	if limit.Limit() != nil && *limit.Limit() > 0 {
		stmt = stmt.LIMIT(*limit.Limit())
	} else {
		stmt = stmt.LIMIT(ListDefaultLimit)
	}

	return stmt
}

type offsetter interface {
	Offset() *int64
}

// setOffset sets an offset to the returned rows in a select statement.
func setOffset(stmt SelectStatement, offset offsetter) SelectStatement {
	if offset.Offset() != nil && *offset.Offset() > 0 {
		stmt = stmt.OFFSET(*offset.Offset())
	}

	return stmt
}
