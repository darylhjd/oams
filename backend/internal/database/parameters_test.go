package database

import (
	"testing"

	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	"github.com/darylhjd/oams/backend/pkg/to"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/stretchr/testify/assert"
)

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

			stmt := setSorts(SELECT(NULL), tt.withSorter)
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

			stmt := setLimit(SELECT(NULL), tt.withLimiter)
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

			stmt := setOffset(SELECT(NULL), tt.withOffset)
			a.Equal(tt.wantStatement.DebugSql(), stmt.DebugSql())
		})
	}
}
