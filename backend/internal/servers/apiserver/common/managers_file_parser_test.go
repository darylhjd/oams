package common

import (
	"os"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/stretchr/testify/assert"
)

func TestParseManagersFile(t *testing.T) {
	tests := []struct {
		name         string
		file         string
		wantErr      string
		expectedData []database.UpsertClassGroupManagerParams
	}{
		{
			"sample manager file",
			"managers_file_well_formatted.xlsx",
			"",
			[]database.UpsertClassGroupManagerParams{
				{
					UserID:         "HARJ0002",
					ClassCode:      "SC1005",
					ClassYear:      2023,
					ClassSemester:  "1",
					ClassGroupName: "A21",
					ClassType:      model.ClassType_Lec,
					ManagingRole:   model.ManagingRole_CourseCoordinator,
				},
				{
					UserID:         "AC2003",
					ClassCode:      "SC4502",
					ClassYear:      2022,
					ClassSemester:  "2",
					ClassGroupName: "E33",
					ClassType:      model.ClassType_Lab,
					ManagingRole:   model.ManagingRole_CourseCoordinator,
				},
			},
		},
		{
			"row has wrong number of columns",
			"managers_file_wrong_col_length.xlsx",
			"unexpected number of data columns",
			nil,
		},
		{
			"header has wrong name",
			"managers_file_wrong_header_name.xlsx",
			"failed sanity check for header name",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			file, err := os.Open(tt.file)
			a.Nil(err)

			data, err := ParseManagersFile(file)
			if tt.wantErr != "" {
				a.Contains(err.Error(), tt.wantErr)
				return
			}

			a.Equal(tt.expectedData, data)
		})
	}
}
