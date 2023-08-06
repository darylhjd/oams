package common

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/database"
)

func TestParseBatchFile(t *testing.T) {
	loc, err := time.LoadLocation(timezoneLocation)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		file         string
		wantErr      string
		expectedData *BatchData
	}{
		{
			"sample class creation file",
			"batch_file_well_formatted.xlsx",
			"",
			&BatchData{
				"batch_file_well_formatted.xlsx",
				time.Date(2023, time.June, 15, 13, 1, 0, 0, loc),
				database.UpsertClassesParams{
					Code:      "SC1015",
					Year:      2022,
					Semester:  "2",
					Programme: "CSC  Full-Time",
					Au:        3,
				},
				[]ClassGroupData{
					{
						database.UpsertClassGroupsParams{
							Name:      "A21",
							ClassType: database.ClassTypeTUT,
						},
						[]SessionData{},
						[]database.UpsertUsersParams{
							{"CHUL6789", "CHUA LI TING", "", database.UserRoleSTUDENT},
							{"YAPW9087", "YAP WEN LI", "", database.UserRoleSTUDENT},
						},
					},
					{
						database.UpsertClassGroupsParams{
							Name:      "A26",
							ClassType: database.ClassTypeTUT,
						},
						[]SessionData{},
						[]database.UpsertUsersParams{
							{"BENST129", "BENJAMIN SANTOS", "", database.UserRoleSTUDENT},
							{"YAPW9088", "YAP WEI LING", "", database.UserRoleSTUDENT},
						},
					},
					{
						database.UpsertClassGroupsParams{
							Name:      "A32",
							ClassType: database.ClassTypeTUT,
						},
						[]SessionData{},
						[]database.UpsertUsersParams{
							{"PATELAR14", "ARJUN PATEL", "", database.UserRoleSTUDENT},
							{"YAPX9087", "YAP XIN TING", "", database.UserRoleSTUDENT},
						},
					},
				},
				database.ClassTypeTUT,
			},
		},
		{
			"empty class creation file",
			"batch_file_empty_test.xlsx",
			"not enough rows for class metadata",
			nil,
		},
		{
			"too many class metadata rows",
			"batch_file_excessive_class_metadata_rows.xlsx",
			"unexpected number of columns for class group row",
			nil,
		},
		{
			"missing class group enrollment list identifier",
			"batch_file_no_enrollment_list_ident.xlsx",
			"unexpected start of class group enrollment list",
			nil,
		},
		{
			"second class group missing enrollment list identifier",
			"batch_file_no_enrollment_list_ident_2.xlsx",
			"unexpected start of class group enrollment list",
			nil,
		},
		{
			"student row with wrong length",
			"batch_file_student_row_wrong_length.xlsx",
			"unexpected number of columns for student enrollment row",
			nil,
		},
		{
			"invalid format for class group name",
			"batch_file_invalid_class_group_name.xlsx",
			"could not parse class group",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			file, err := os.Open(tt.file)
			a.Nil(err)

			data, err := ParseBatchFile(tt.file, file)
			if tt.wantErr != "" {
				a.Contains(err.Error(), tt.wantErr)
				return
			}

			attributeTests := map[any]any{
				tt.expectedData.Filename:         data.Filename,
				tt.expectedData.FileCreationDate: data.FileCreationDate,
				tt.expectedData.Class:            data.Class,
				tt.expectedData.classType:        data.classType,
			}

			for expected, actual := range attributeTests {
				a.Equal(expected, actual)
			}

			// Check class groups.
			a.Equal(len(tt.expectedData.ClassGroups), len(data.ClassGroups))
			for i := 0; i < len(tt.expectedData.ClassGroups); i++ {
				a.Equal(tt.expectedData.ClassGroups[i].Name, data.ClassGroups[i].Name)
				a.Equal(tt.expectedData.ClassGroups[i].ClassType, data.ClassGroups[i].ClassType)
				// TODO: Test for sessions.
				for _, student := range tt.expectedData.ClassGroups[i].Students {
					a.Contains(data.ClassGroups[i].Students, student)
				}
			}
		})
	}
}
