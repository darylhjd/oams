package database

import (
	"fmt"
	"strings"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/gorilla/schema"
	"github.com/lib/pq"
)

var decoder = schema.NewDecoder()

type ListQueryParams struct {
	F       []string `schema:"filter"`
	FParsed []BoolExpression
	S       []string `schema:"sort"`
	SParsed []SortParam
	L       *int64 `schema:"limit"`
	O       *int64 `schema:"offset"`
}

type SortParam struct {
	Col          Column
	IsDescending bool
}

// setFilters sets the WHERE parameters for a select statement.
func (p ListQueryParams) setFilters(stmt SelectStatement) SelectStatement {
	if len(p.FParsed) == 0 {
		return stmt
	}

	return stmt.WHERE(AND(p.FParsed...))
}

// setSorts sets the ORDER BY parameters for a select statement.
func (p ListQueryParams) setSorts(stmt SelectStatement) SelectStatement {
	if len(p.SParsed) == 0 {
		return stmt
	}

	var orders []OrderByClause
	for _, param := range p.SParsed {
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
func (p ListQueryParams) setLimit(stmt SelectStatement) SelectStatement {
	if p.L != nil && *p.L > 0 {
		stmt = stmt.LIMIT(*p.L)
	} else {
		stmt = stmt.LIMIT(ListDefaultLimit)
	}

	return stmt
}

// setOffset sets an offset to the returned rows in a select statement.
func (p ListQueryParams) setOffset(stmt SelectStatement) SelectStatement {
	if p.O != nil && *p.O > 0 {
		stmt = stmt.OFFSET(*p.O)
	}

	return stmt
}

func DecodeListQueryParams(source map[string][]string, table Table, cList ColumnList) (ListQueryParams, error) {
	var l ListQueryParams
	err := decoder.Decode(&l, source)
	if err != nil {
		return l, err
	}

	for _, f := range l.F {
		p, err := parseFilterParam(f, table, cList)
		if err != nil {
			return l, err
		}

		l.FParsed = append(l.FParsed, p)
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
	filterSplitLength = 2
	filterSeparator   = "."
)

func parseFilterParam(s string, table Table, cList ColumnList) (BoolExpression, error) {
	fmt.Println(s)
	parts := strings.Split(s, filterSeparator)
	if len(parts) != filterSplitLength {
		return nil, fmt.Errorf("filter param `%s` does not follow format `col_name.value`", s)
	}

	for _, c := range cList {
		if c.Name() == parts[0] {
			return RawBool(fmt.Sprintf("%s.%s = #arg", pq.QuoteIdentifier(table.Alias()), c.Name()), RawArgs{
				"#arg": parts[1],
			}), nil
		}
	}

	return nil, fmt.Errorf("invalid column `%s`", parts[0])
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
