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
		containsErr  string
		expectedData ClassCreationData
	}{
		{
			"sample class lab creation file",
			"class_lab_test.xlsx",
			"",
			ClassCreationData{
				2022,
				"2",
				time.Date(2023, time.June, 15, 13, 1, 0, 0, loc),
				"SC1015",
				"LAB",
				[]ClassGroup{
					{
						"A21",
						[]database.ClassVenueSchedule{},
						[]database.Student{
							{"CHUL6789", "CHUA LI TING", pgtype.Text{}},
							{"YAPW9087", "YAP WEN LI", pgtype.Text{}},
						},
					},
					{
						"A26",
						[]database.ClassVenueSchedule{},
						[]database.Student{
							{"BENST129", "BENJAMIN SANTOS", pgtype.Text{}},
							{"YAPW9087", "YAP WEI LING", pgtype.Text{}},
						},
					},
					{
						"A32",
						[]database.ClassVenueSchedule{},
						[]database.Student{
							{"PATELAR14", "ARJUN PATEL", pgtype.Text{}},
							{"YAPX9087", "YAP XIN TING", pgtype.Text{}},
						},
					},
				},
			},
		},
		{
			"sample class lecture creation file",
			"class_lec_test.xlsx",
			"",
			ClassCreationData{
				2022,
				"2",
				time.Date(2023, time.June, 15, 13, 1, 0, 0, loc),
				"SC1015",
				"LEC",
				[]ClassGroup{
					{
						"L1",
						[]database.ClassVenueSchedule{},
						[]database.Student{
							{"PATELAR14", "ARJUN PATEL", pgtype.Text{}},
							{"YAPX9087", "YAP XIN TING", pgtype.Text{}},
						},
					},
				},
			},
		},
		{
			"sample class tutorial creation file",
			"class_tut_test.xlsx",
			"",
			ClassCreationData{
				2022,
				"2",
				time.Date(2023, time.June, 15, 13, 1, 0, 0, loc),
				"SC1015",
				"TUT",
				[]ClassGroup{
					{
						"A21",
						[]database.ClassVenueSchedule{},
						[]database.Student{
							{"CHUL6789", "CHUA LI TING", pgtype.Text{}},
							{"YAPW9087", "YAP WEN LI", pgtype.Text{}},
						},
					},
					{
						"A26",
						[]database.ClassVenueSchedule{},
						[]database.Student{
							{"BENST129", "BENJAMIN SANTOS", pgtype.Text{}},
							{"YAPW9087", "YAP WEI LING", pgtype.Text{}},
						},
					},
					{
						"A32",
						[]database.ClassVenueSchedule{},
						[]database.Student{
							{"PATELAR14", "ARJUN PATEL", pgtype.Text{}},
							{"YAPX9087", "YAP XIN TING", pgtype.Text{}},
						},
					},
				},
			},
		},
		{
			"empty class creation file",
			"class_empty_test.xlsx",
			"not enough rows for class metadata",
			ClassCreationData{},
		},
		{
			"too many class metadata rows",
			"class_excessive_class_metadata_rows.xlsx",
			"unexpected number of columns for class group row",
			ClassCreationData{},
		},
		{
			"missing class group enrollment list identifier",
			"class_missing_enrollment_ident.xlsx",
			"unexpected start of class group enrollment list",
			ClassCreationData{},
		},
		{
			"second class group missing enrollment list identifier",
			"class_second_group_missing_enrollment_ident.xlsx",
			"unexpected start of class group enrollment list",
			ClassCreationData{},
		},
		{
			"student row with wrong length",
			"class_student_row_wrong_length.xlsx",
			"unexpected number of columns for student enrollment row",
			ClassCreationData{},
		},
		{
			"class group with no enrollment",
			"class_group_with_no_enrollment.xlsx",
			"class group A21 has no enrollments",
			ClassCreationData{},
		},
		{
			"course with no class groups",
			"class_with_no_groups.xlsx",
			"class creation file has no valid class groups",
			ClassCreationData{},
		},
		{
			"invalid format for class group name",
			"class_with_invalid_class_group_name.xlsx",
			"could not parse class group",
			ClassCreationData{},
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.file)
			a.Nil(err)

			data, err := ParseClassCreationFile(file)
			if tt.containsErr != "" {
				a.ErrorContains(err, tt.containsErr)
				return
			}

			attributeTests := map[any]any{
				tt.expectedData.Year:             data.Year,
				tt.expectedData.Semester:         data.Semester,
				tt.expectedData.FileCreationDate: data.FileCreationDate,
				tt.expectedData.CourseCode:       data.CourseCode,
				tt.expectedData.ClassType:        data.ClassType,
			}

			for expected, actual := range attributeTests {
				a.Equal(expected, actual)
			}

			// Check class groups.
			a.Equal(len(tt.expectedData.ClassGroups), len(data.ClassGroups))
			for i := 0; i < len(tt.expectedData.ClassGroups); i++ {
				a.Equal(tt.expectedData.ClassGroups[i].Name, data.ClassGroups[i].Name)
				// TODO: Test for schedule.
				for _, student := range tt.expectedData.ClassGroups[i].Students {
					a.Contains(data.ClassGroups[i].Students, student)
				}
			}
		})
	}
}
