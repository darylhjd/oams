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
		wantErr      string
		expectedData *ClassCreationData
	}{
		{
			"sample class lab creation file",
			"class_lab_test.xlsx",
			"",
			&ClassCreationData{
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
						database.UpsertClassGroupsParams{
							Name:      "A21",
							ClassType: database.ClassTypeLAB,
						},
						[]SessionData{},
						[]database.UpsertStudentsParams{
							{"CHUL6789", "CHUA LI TING", pgtype.Text{}},
							{"YAPW9087", "YAP WEN LI", pgtype.Text{}},
						},
					},
					{
						database.UpsertClassGroupsParams{
							Name:      "A26",
							ClassType: database.ClassTypeLAB,
						},
						[]SessionData{},
						[]database.UpsertStudentsParams{
							{"BENST129", "BENJAMIN SANTOS", pgtype.Text{}},
							{"YAPW9087", "YAP WEI LING", pgtype.Text{}},
						},
					},
					{
						database.UpsertClassGroupsParams{
							Name:      "A32",
							ClassType: database.ClassTypeLAB,
						},
						[]SessionData{},
						[]database.UpsertStudentsParams{
							{"PATELAR14", "ARJUN PATEL", pgtype.Text{}},
							{"YAPX9087", "YAP XIN TING", pgtype.Text{}},
						},
					},
				},
				database.ClassTypeLAB,
			},
		},
		{
			"sample class lecture creation file",
			"class_lec_test.xlsx",
			"",
			&ClassCreationData{
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
						database.UpsertClassGroupsParams{
							Name:      "L1",
							ClassType: database.ClassTypeLEC,
						},
						[]SessionData{},
						[]database.UpsertStudentsParams{
							{"PATELAR14", "ARJUN PATEL", pgtype.Text{}},
							{"YAPX9087", "YAP XIN TING", pgtype.Text{}},
						},
					},
				},
				database.ClassTypeLEC,
			},
		},
		{
			"sample class tutorial creation file",
			"class_tut_test.xlsx",
			"",
			&ClassCreationData{
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
						database.UpsertClassGroupsParams{
							Name:      "A21",
							ClassType: database.ClassTypeTUT,
						},
						[]SessionData{},
						[]database.UpsertStudentsParams{
							{"CHUL6789", "CHUA LI TING", pgtype.Text{}},
							{"YAPW9087", "YAP WEN LI", pgtype.Text{}},
						},
					},
					{
						database.UpsertClassGroupsParams{
							Name:      "A26",
							ClassType: database.ClassTypeTUT,
						},
						[]SessionData{},
						[]database.UpsertStudentsParams{
							{"BENST129", "BENJAMIN SANTOS", pgtype.Text{}},
							{"YAPW9087", "YAP WEI LING", pgtype.Text{}},
						},
					},
					{
						database.UpsertClassGroupsParams{
							Name:      "A32",
							ClassType: database.ClassTypeTUT,
						},
						[]SessionData{},
						[]database.UpsertStudentsParams{
							{"PATELAR14", "ARJUN PATEL", pgtype.Text{}},
							{"YAPX9087", "YAP XIN TING", pgtype.Text{}},
						},
					},
				},
				database.ClassTypeTUT,
			},
		},
		{
			"empty class creation file",
			"class_empty_test.xlsx",
			"not enough rows for class metadata",
			nil,
		},
		{
			"too many class metadata rows",
			"class_excessive_class_metadata_rows.xlsx",
			"unexpected number of columns for class group row",
			nil,
		},
		{
			"missing class group enrollment list identifier",
			"class_missing_enrollment_ident.xlsx",
			"unexpected start of class group enrollment list",
			nil,
		},
		{
			"second class group missing enrollment list identifier",
			"class_second_group_missing_enrollment_ident.xlsx",
			"unexpected start of class group enrollment list",
			nil,
		},
		{
			"student row with wrong length",
			"class_student_row_wrong_length.xlsx",
			"unexpected number of columns for student enrollment row",
			nil,
		},
		{
			"class group with no enrollment",
			"class_group_with_no_enrollment.xlsx",
			"class group A21 has no enrollments",
			nil,
		},
		{
			"course with no class groups",
			"class_with_no_groups.xlsx",
			"creation data has no valid class groups",
			nil,
		},
		{
			"invalid format for class group name",
			"class_with_invalid_class_group_name.xlsx",
			"could not parse class group",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			file, err := os.Open(tt.file)
			a.Nil(err)

			data, err := ParseClassCreationFile(tt.file, file)
			if tt.wantErr != "" {
				a.Contains(err.Error(), tt.wantErr)
				return
			}

			attributeTests := map[any]any{
				tt.expectedData.Filename:         data.Filename,
				tt.expectedData.FileCreationDate: data.FileCreationDate,
				tt.expectedData.Course:           data.Course,
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
