package common

import (
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/database"
)

func TestParseClassCreationFile(t *testing.T) {
	loc, err := time.LoadLocation(timezoneLocation)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name         string
		file         string
		expectedData ClassCreationData
	}{
		{
			"sample class lab creation file",
			"class_lab_test.xlsx",
			ClassCreationData{
				"class_lab_test.xlsx",
				time.Date(2023, time.June, 15, 13, 1, 0, 0, loc),
				database.UpsertCoursesParams{
					Code:      "SC1015",
					Year:      2022,
					Semester:  "2",
					Programme: "CSC  Full-Time",
					Au:        3,
				},
				[]ClassGroupData{
					{
						database.CreateClassGroupsParams{
							Name:      "A21",
							ClassType: database.ClassTypeLAB,
						},
						[]database.CreateClassGroupSessionsParams{},
						[]database.UpsertStudentsParams{
							{"CHUL6789", "CHUA LI TING", pgtype.Text{}},
							{"YAPW9087", "YAP WEN LI", pgtype.Text{}},
						},
					},
					{
						database.CreateClassGroupsParams{
							Name:      "A26",
							ClassType: database.ClassTypeLAB,
						},
						[]database.CreateClassGroupSessionsParams{},
						[]database.UpsertStudentsParams{
							{"BENST129", "BENJAMIN SANTOS", pgtype.Text{}},
							{"YAPW9087", "YAP WEI LING", pgtype.Text{}},
						},
					},
					{
						database.CreateClassGroupsParams{
							Name:      "A32",
							ClassType: database.ClassTypeLAB,
						},
						[]database.CreateClassGroupSessionsParams{},
						[]database.UpsertStudentsParams{
							{"PATELAR14", "ARJUN PATEL", pgtype.Text{}},
							{"YAPX9087", "YAP XIN TING", pgtype.Text{}},
						},
					},
				},
				"",
				database.ClassTypeLAB,
			},
		},
		{
			"sample class lecture creation file",
			"class_lec_test.xlsx",
			ClassCreationData{
				"class_lec_test.xlsx",
				time.Date(2023, time.June, 15, 13, 1, 0, 0, loc),
				database.UpsertCoursesParams{
					Code:      "SC1015",
					Year:      2022,
					Semester:  "2",
					Programme: "CSC  Full-Time",
					Au:        3,
				},
				[]ClassGroupData{
					{
						database.CreateClassGroupsParams{
							Name:      "L1",
							ClassType: database.ClassTypeLEC,
						},
						[]database.CreateClassGroupSessionsParams{},
						[]database.UpsertStudentsParams{
							{"PATELAR14", "ARJUN PATEL", pgtype.Text{}},
							{"YAPX9087", "YAP XIN TING", pgtype.Text{}},
						},
					},
				},
				"",
				database.ClassTypeLEC,
			},
		},
		{
			"sample class tutorial creation file",
			"class_tut_test.xlsx",
			ClassCreationData{
				"class_tut_test.xlsx",
				time.Date(2023, time.June, 15, 13, 1, 0, 0, loc),
				database.UpsertCoursesParams{
					Code:      "SC1015",
					Year:      2022,
					Semester:  "2",
					Programme: "CSC  Full-Time",
					Au:        3,
				},
				[]ClassGroupData{
					{
						database.CreateClassGroupsParams{
							Name:      "A21",
							ClassType: database.ClassTypeTUT,
						},
						[]database.CreateClassGroupSessionsParams{},
						[]database.UpsertStudentsParams{
							{"CHUL6789", "CHUA LI TING", pgtype.Text{}},
							{"YAPW9087", "YAP WEN LI", pgtype.Text{}},
						},
					},
					{
						database.CreateClassGroupsParams{
							Name:      "A26",
							ClassType: database.ClassTypeTUT,
						},
						[]database.CreateClassGroupSessionsParams{},
						[]database.UpsertStudentsParams{
							{"BENST129", "BENJAMIN SANTOS", pgtype.Text{}},
							{"YAPW9087", "YAP WEI LING", pgtype.Text{}},
						},
					},
					{
						database.CreateClassGroupsParams{
							Name:      "A32",
							ClassType: database.ClassTypeTUT,
						},
						[]database.CreateClassGroupSessionsParams{},
						[]database.UpsertStudentsParams{
							{"PATELAR14", "ARJUN PATEL", pgtype.Text{}},
							{"YAPX9087", "YAP XIN TING", pgtype.Text{}},
						},
					},
				},
				"",
				database.ClassTypeTUT,
			},
		},
		{
			"empty class creation file",
			"class_empty_test.xlsx",
			ClassCreationData{
				Error: "not enough rows for class metadata",
			},
		},
		{
			"too many class metadata rows",
			"class_excessive_class_metadata_rows.xlsx",
			ClassCreationData{
				Error: "unexpected number of columns for class group row",
			},
		},
		{
			"missing class group enrollment list identifier",
			"class_missing_enrollment_ident.xlsx",
			ClassCreationData{
				Error: "unexpected start of class group enrollment list",
			},
		},
		{
			"second class group missing enrollment list identifier",
			"class_second_group_missing_enrollment_ident.xlsx",
			ClassCreationData{
				Error: "unexpected start of class group enrollment list",
			},
		},
		{
			"student row with wrong length",
			"class_student_row_wrong_length.xlsx",
			ClassCreationData{
				Error: "unexpected number of columns for student enrollment row",
			},
		},
		{
			"class group with no enrollment",
			"class_group_with_no_enrollment.xlsx",
			ClassCreationData{
				Error: "class group A21 has no enrollments",
			},
		},
		{
			"course with no class groups",
			"class_with_no_groups.xlsx",
			ClassCreationData{
				Error: "creation data has no valid class groups",
			},
		},
		{
			"invalid format for class group name",
			"class_with_invalid_class_group_name.xlsx",
			ClassCreationData{
				Error: "could not parse class group",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			file, err := os.Open(tt.file)
			a.Nil(err)

			data, err := ParseClassCreationFile(tt.file, file)
			a.Nil(err)

			if tt.expectedData.Error != "" {
				a.Contains(data.Error, tt.expectedData.Error)
				return
			}

			attributeTests := map[any]any{
				tt.expectedData.Filename:         data.Filename,
				tt.expectedData.FileCreationDate: data.FileCreationDate,
				tt.expectedData.Course:           data.Course,
				tt.expectedData.ClassType:        data.ClassType,
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
