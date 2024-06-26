package database

import (
	"testing"

	. "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
	"github.com/darylhjd/oams/backend/pkg/to"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/stretchr/testify/assert"
)

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
			SELECT(NULL),
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
