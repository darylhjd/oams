package v1

import (
	"fmt"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/go-jet/jet/v2/postgres"
)

// listQueryParameters is a generic struct that stores the query parameters for all GET list endpoints.
type listQueryParameters struct {
	S       []string `schema:"sort"`
	SParsed []database.SortParam
	L       *int64 `schema:"limit"`
	O       *int64 `schema:"offset"`
}

func (p listQueryParameters) Sorts() []database.SortParam {
	return p.SParsed
}

func (p listQueryParameters) Limit() *int64 {
	return p.L
}

func (p listQueryParameters) Offset() *int64 {
	return p.O
}

const (
	sortSeparator = ":"
)

func (v *APIServerV1) decodeListQueryParameters(src map[string][]string, cList postgres.ColumnList) (listQueryParameters, error) {
	var l listQueryParameters
	err := v.decoder.Decode(&l, src)
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
	maxSortParamLength = 2
)

// parseSortParam is a helper function to get a database.SortParam from one sort parameter.
func parseSortParam(s string, cList postgres.ColumnList) (database.SortParam, error) {
	var p database.SortParam

	sort := strings.Split(s, sortSeparator)
	if len(sort) > maxSortParamLength || (len(sort) == maxSortParamLength && sort[maxSortParamLength-1] != database.SortDirectionDesc) {
		return p, fmt.Errorf(
			"sort parameter should be of format `col` (default asc) or `col%s[%s]",
			sortSeparator,
			database.SortDirectionDesc,
		)
	}

	var found bool
	for _, c := range cList {
		if c.Name() == sort[0] {
			found = true
			p.Col = c
		}
	}

	switch {
	case !found:
		return p, fmt.Errorf("unknown sort column `%s`", sort[0])
	case len(sort) == maxSortParamLength:
		p.Direction = sort[maxSortParamLength-1]
		fallthrough
	default:
		return p, nil
	}
}
