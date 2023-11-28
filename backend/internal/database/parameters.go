package database

import (
	"fmt"
	"strings"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type ListQueryParams struct {
	S       []string `schema:"sort"`
	SParsed []SortParam
	L       *int64   `schema:"limit"`
	O       *int64   `schema:"offset"`
	F       []string `schema:"filter"`
}

type SortParam struct {
	Col          Column
	IsDescending bool
}

func DecodeListQueryParams(source map[string][]string, cList ColumnList) (ListQueryParams, error) {
	var l ListQueryParams
	err := decoder.Decode(&l, source)
	if err != nil {
		return l, err
	}

	for _, s := range l.S {
		p, err := parseSortParam(s, cList)
		if err != nil {
			return l, err
		}

		l.SParsed = append(l.SParsed, p)
	}

	return l, err
}

const (
	descendingIdent = ":desc"
)

func parseSortParam(s string, cList ColumnList) (SortParam, error) {
	var p SortParam

	if strings.HasSuffix(s, descendingIdent) {
		p.IsDescending = true
	}

	columnString := strings.TrimSuffix(s, descendingIdent)
	for _, c := range cList {
		if c.Name() == columnString {
			p.Col = c
			return p, nil
		}
	}

	return p, fmt.Errorf("unknown sort column `%s`", columnString)
}

// setSorts sets the ORDER BY parameters for a select statement.
func setSorts(stmt SelectStatement, params ListQueryParams) SelectStatement {
	if len(params.SParsed) == 0 {
		return stmt
	}

	var orders []OrderByClause
	for _, param := range params.SParsed {
		direction := param.Col.ASC()
		if param.IsDescending {
			direction = param.Col.DESC()
		}

		orders = append(orders, direction)
	}

	return stmt.ORDER_BY(orders...)
}

const (
	ListDefaultLimit = 50
)

// setLimit sets a limit to the number of returned rows in a select statement.
func setLimit(stmt SelectStatement, params ListQueryParams) SelectStatement {
	if params.L != nil && *params.L > 0 {
		stmt = stmt.LIMIT(*params.L)
	} else {
		stmt = stmt.LIMIT(ListDefaultLimit)
	}

	return stmt
}

// setOffset sets an offset to the returned rows in a select statement.
func setOffset(stmt SelectStatement, params ListQueryParams) SelectStatement {
	if params.O != nil && *params.O > 0 {
		stmt = stmt.OFFSET(*params.O)
	}

	return stmt
}
