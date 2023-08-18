package database

import (
	. "github.com/go-jet/jet/v2/postgres"
)

type listParams interface {
	sorter
	limiter
	offsetter
}

const (
	SortDirectionDesc = "desc"
)

type sorter interface {
	Sorts() []SortParam
}

type SortParam struct {
	Col       Column
	Direction string
}

// setSorts sets the ORDER BY parameters for a select statement.
func setSorts(stmt SelectStatement, s sorter) SelectStatement {
	if len(s.Sorts()) == 0 {
		return stmt
	}

	var orders []OrderByClause
	for _, param := range s.Sorts() {
		direction := param.Col.ASC()
		if param.Direction == SortDirectionDesc {
			direction = param.Col.DESC()
		}

		orders = append(orders, direction)
	}

	return stmt.ORDER_BY(orders...)
}

const (
	ListDefaultLimit = 50
)

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
