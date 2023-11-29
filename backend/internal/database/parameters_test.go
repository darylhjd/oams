package database

import (
	"fmt"
	"testing"

	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	"github.com/darylhjd/oams/backend/pkg/to"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_parseFilterParam(t *testing.T) {
	tts := []struct {
		name            string
		withF           string
		withTable       Table
		withColumnList  ColumnList
		wantFilterParam string
		wantErr         bool
	}{
		{
			"correct column",
			"id=NTU0001",
			Users,
			Users.AllColumns,
			fmt.Sprintf("%s.id = 'NTU0001'", pq.QuoteIdentifier(Users.Alias())),
			false,
		},
		{
			"correct column integer",
			"id=2",
			Classes,
			Classes.AllColumns,
			fmt.Sprintf("%s.id = '2'", pq.QuoteIdentifier(Classes.Alias())),
			false,
		},
		{
			"wrong column",
			"wrong=NTU0001",
			Users,
			Users.AllColumns,
			"",
			true,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			filterParam, err := parseFilterParam(tt.withF, tt.withTable, tt.withColumnList)
			if tt.wantErr {
				a.Error(err)
			} else {
				a.Nil(err)
				a.Contains(
					SELECT(NULL).WHERE(filterParam).DebugSql(),
					tt.wantFilterParam,
				)
			}
		})
	}
}

func Test_parseSortParam(t *testing.T) {
	tts := []struct {
		name           string
		withS          string
		withColumnList ColumnList
		wantSortParam  SortParam
		wantErr        bool
	}{
		{
			"correct column descending",
			"id:desc",
			Users.AllColumns,
			SortParam{
				Users.ID,
				true,
			},
			false,
		},
		{
			"correct column ascending",
			"id",
			Users.AllColumns,
			SortParam{
				Users.ID,
				false,
			},
			false,
		},
		{
			"wrong column",
			"wrong",
			Users.AllColumns,
			SortParam{},
			true,
		},
		{
			"wrong descending identifier",
			"id:asc",
			Users.AllColumns,
			SortParam{},
			true,
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			sortParam, err := parseSortParam(tt.withS, tt.withColumnList)
			if tt.wantErr {
				a.Error(err)
			} else {
				a.Nil(err)
				a.Equal(tt.wantSortParam, sortParam)
			}
		})
	}
}

func Test_setFilters(t *testing.T) {
	tts := []struct {
		name          string
		withFilter    ListQueryParams
		wantStatement SelectStatement
	}{
		{
			"with parameters",
			ListQueryParams{
				FParsed: []BoolExpression{
					RawBool(fmt.Sprintf("%s.id = 'NTU0001'", pq.QuoteIdentifier(Users.Alias()))),
					RawBool(fmt.Sprintf("%s.id = 'NTU0002'", pq.QuoteIdentifier(Users.Alias()))),
				},
			},
			SELECT(NULL).WHERE(
				AND(
					RawBool(fmt.Sprintf("%s.id = 'NTU0001'", pq.QuoteIdentifier(Users.Alias()))),
					RawBool(fmt.Sprintf("%s.id = 'NTU0002'", pq.QuoteIdentifier(Users.Alias()))),
				),
			),
		},
		{
			"with no parameters",
			ListQueryParams{
				FParsed: []BoolExpression{},
			},
			SELECT(NULL),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			stmt := tt.withFilter.setFilters(SELECT(NULL))
			a.Equal(tt.wantStatement.DebugSql(), stmt.DebugSql())
		})
	}
}

func Test_setSorts(t *testing.T) {
	tts := []struct {
		name          string
		withSorter    ListQueryParams
		wantStatement SelectStatement
	}{
		{
			"with parameters",
			ListQueryParams{
				SParsed: []SortParam{
					{Users.ID, false},
					{Users.Name, true},
				}},
			SELECT(NULL).
				ORDER_BY(
					Users.ID.ASC(),
					Users.Name.DESC(),
				),
		},
		{
			"with no parameters",
			ListQueryParams{},
			SELECT(NULL),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			stmt := tt.withSorter.setSorts(SELECT(NULL))
			a.Equal(tt.wantStatement.DebugSql(), stmt.DebugSql())
		})
	}
}

func Test_setLimit(t *testing.T) {
	tts := []struct {
		name          string
		withLimiter   ListQueryParams
		wantStatement SelectStatement
	}{
		{
			"with parameters",
			ListQueryParams{L: to.Ptr(int64(2))},
			SELECT(NULL).
				LIMIT(2),
		},
		{
			"with no parameters",
			ListQueryParams{},
			SELECT(NULL).
				LIMIT(ListDefaultLimit),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			stmt := tt.withLimiter.setLimit(SELECT(NULL))
			a.Equal(tt.wantStatement.DebugSql(), stmt.DebugSql())
		})
	}
}

func Test_setOffset(t *testing.T) {
	tts := []struct {
		name          string
		withOffset    ListQueryParams
		wantStatement SelectStatement
	}{
		{
			"with parameters",
			ListQueryParams{O: to.Ptr(int64(2))},
			SELECT(NULL).
				OFFSET(2),
		},
		{
			"with no parameters",
			ListQueryParams{},
			SELECT(NULL),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			stmt := tt.withOffset.setOffset(SELECT(NULL))
			a.Equal(tt.wantStatement.DebugSql(), stmt.DebugSql())
		})
	}
}
