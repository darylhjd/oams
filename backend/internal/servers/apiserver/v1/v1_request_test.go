package v1

import (
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/stretchr/testify/assert"
)

func Test_parseSortParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		withString  string
		withColList postgres.ColumnList
		wantParam   database.SortParam
		wantErr     string
	}{
		{
			"correct column name ascending",
			"name",
			table.Users.AllColumns,
			database.SortParam{
				Col:       table.Users.Name,
				Direction: "",
			},
			"",
		},
		{
			"correct column name descending",
			"id:desc",
			table.Users.AllColumns,
			database.SortParam{
				Col:       table.Users.ID,
				Direction: "desc",
			},
			"",
		},
		{
			"wrong column name",
			"wrong",
			table.Users.AllColumns,
			database.SortParam{},
			"unknown sort column `wrong`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			param, err := parseSortParam(tt.withString, tt.withColList)
			if tt.wantErr != "" {
				a.ErrorContains(err, tt.wantErr)
			} else {
				a.Equal(tt.wantParam, param)
			}
		})
	}
}
