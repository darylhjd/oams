package common

import (
	"os"
	"testing"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/database"
)

func TestParseBatchFile(t *testing.T) {
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
				time.Date(2023, time.June, 15, 13, 1, 0, 0, location),
				database.UpsertClassParams{
					Code:      "SC1015",
					Year:      2022,
					Semester:  "2",
					Programme: "CSC  Full-Time",
					Au:        3,
				},
				[]ClassGroupData{
					{
						database.UpsertClassGroupParams{
							Name:      "A21",
							ClassType: model.ClassType_Tut,
						},
						[]database.UpsertClassGroupSessionParams{
							{
								StartTime: time.Date(2023, time.January, 16, 8, 30, 0, 0, location),
								EndTime:   time.Date(2023, time.January, 16, 9, 20, 0, 0, location),
								Venue:     "TR+15 NORTH,NS4-05-93",
							},
							{
								StartTime: time.Date(2023, time.April, 10, 8, 30, 0, 0, location),
								EndTime:   time.Date(2023, time.April, 10, 9, 20, 0, 0, location),
								Venue:     "TR+15 NORTH,NS4-05-93",
							},
							{
								StartTime: time.Date(2023, time.January, 17, 9, 30, 0, 0, location),
								EndTime:   time.Date(2023, time.January, 17, 10, 20, 0, 0, location),
								Venue:     "TR+15 NORTH,NS4-05-93",
							},
							{
								StartTime: time.Date(2023, time.April, 11, 9, 30, 0, 0, location),
								EndTime:   time.Date(2023, time.April, 11, 10, 20, 0, 0, location),
								Venue:     "TR+15 NORTH,NS4-05-93",
							},
						},
						[]database.UpsertUserParams{
							{"CHUL6789", "CHUA LI TING", ""},
							{"YAPW9087", "YAP WEN LI", ""},
						},
					},
					{
						database.UpsertClassGroupParams{
							Name:      "A26",
							ClassType: model.ClassType_Tut,
						},
						[]database.UpsertClassGroupSessionParams{
							{
								StartTime: time.Date(2023, time.January, 16, 1, 30, 0, 0, location),
								EndTime:   time.Date(2023, time.January, 16, 2, 20, 0, 0, location),
								Venue:     "TR+19 NORTH,NS4-05-97",
							},
							{
								StartTime: time.Date(2023, time.April, 10, 1, 30, 0, 0, location),
								EndTime:   time.Date(2023, time.April, 10, 2, 20, 0, 0, location),
								Venue:     "TR+19 NORTH,NS4-05-97",
							},
						},
						[]database.UpsertUserParams{
							{"BENST129", "BENJAMIN SANTOS", ""},
							{"YAPW9088", "YAP WEI LING", ""},
						},
					},
					{
						database.UpsertClassGroupParams{
							Name:      "A32",
							ClassType: model.ClassType_Tut,
						},
						[]database.UpsertClassGroupSessionParams{
							{
								StartTime: time.Date(2023, time.January, 17, 11, 30, 0, 0, location),
								EndTime:   time.Date(2023, time.January, 17, 12, 20, 0, 0, location),
								Venue:     "TR+37 NORTH,NS2-05-30",
							},
							{
								StartTime: time.Date(2023, time.January, 31, 11, 30, 0, 0, location),
								EndTime:   time.Date(2023, time.January, 31, 12, 20, 0, 0, location),
								Venue:     "TR+37 NORTH,NS2-05-30",
							},
							{
								StartTime: time.Date(2023, time.March, 21, 11, 30, 0, 0, location),
								EndTime:   time.Date(2023, time.March, 21, 12, 20, 0, 0, location),
								Venue:     "TR+37 NORTH,NS2-05-30",
							},
							{
								StartTime: time.Date(2023, time.April, 4, 11, 30, 0, 0, location),
								EndTime:   time.Date(2023, time.April, 4, 12, 20, 0, 0, location),
								Venue:     "TR+37 NORTH,NS2-05-30",
							},
						},
						[]database.UpsertUserParams{
							{"PATELAR14", "ARJUN PATEL", ""},
							{"YAPX9087", "YAP XIN TING", ""},
						},
					},
				},
				model.ClassType_Tut,
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
				for _, session := range tt.expectedData.ClassGroups[i].Sessions {
					a.Contains(data.ClassGroups[i].Sessions, session)
				}
				for _, student := range tt.expectedData.ClassGroups[i].Students {
					a.Contains(data.ClassGroups[i].Students, student)
				}
			}
		})
	}
}
