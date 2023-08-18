package database

import (
	"testing"

	. "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	"github.com/darylhjd/oams/backend/pkg/to"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/stretchr/testify/assert"
)

type testSorter struct {
	s []SortParam
}

func (t testSorter) Sorts() []SortParam {
	return t.s
}

func Test_setSorts(t *testing.T) {
	tts := []struct {
		name          string
		withSorter    sorter
		wantStatement SelectStatement
	}{
		{
			"with parameters",
			testSorter{[]SortParam{
				{Users.ID, ""},
				{Users.Name, SortDirectionDesc},
			}},
			SELECT(NULL).
				ORDER_BY(
					Users.ID.ASC(),
					Users.Name.DESC(),
				),
		},
		{
			"with no parameters",
			testSorter{[]SortParam{}},
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

type testLimiter struct {
	l *int64
}

func (t testLimiter) Limit() *int64 {
	return t.l
}

func Test_setLimit(t *testing.T) {
	tts := []struct {
		name          string
		withLimiter   limiter
		wantStatement SelectStatement
	}{
		{
			"with parameters",
			testLimiter{to.Ptr(int64(2))},
			SELECT(NULL).
				LIMIT(2),
		},
		{
			"with no parameters",
			testLimiter{},
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

type testOffsetter struct {
	o *int64
}

func (t testOffsetter) Offset() *int64 {
	return t.o
}

func Test_setOffset(t *testing.T) {
	tts := []struct {
		name          string
		withOffsetter offsetter
		wantStatement SelectStatement
	}{
		{
			"with parameters",
			testOffsetter{to.Ptr(int64(2))},
			SELECT(NULL).
				OFFSET(2),
		},
		{
			"with no parameters",
			testOffsetter{},
			SELECT(NULL),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			stmt := setOffset(SELECT(NULL), tt.withOffsetter)
			a.Equal(tt.wantStatement.DebugSql(), stmt.DebugSql())
		})
	}
}
