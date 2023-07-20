package common

import (
	"database/sql"
	"os"
	"testing"
	"time"

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
							{"CHUL6789", "CHUA LI TING", sql.NullString{}},
							{"YAPW9087", "YAP WEN LI", sql.NullString{}},
						},
					},
					{
						"A26",
						[]database.ClassVenueSchedule{},
						[]database.Student{
							{"BENST129", "BENJAMIN SANTOS", sql.NullString{}},
							{"YAPW9087", "YAP WEI LING", sql.NullString{}},
						},
					},
					{
						"A32",
						[]database.ClassVenueSchedule{},
						[]database.Student{
							{"PATELAR14", "ARJUN PATEL", sql.NullString{}},
							{"YAPX9087", "YAP XIN TING", sql.NullString{}},
						},
					},
				},
			},
		},
		{
			"sample class lecture creation file",
			"class_lec_test.xlsx",
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
							{"PATELAR14", "ARJUN PATEL", sql.NullString{}},
							{"YAPX9087", "YAP XIN TING", sql.NullString{}},
						},
					},
				},
			},
		},
		{
			"sample class tutorial creation file",
			"class_tut_test.xlsx",
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
							{"CHUL6789", "CHUA LI TING", sql.NullString{}},
							{"YAPW9087", "YAP WEN LI", sql.NullString{}},
						},
					},
					{
						"A26",
						[]database.ClassVenueSchedule{},
						[]database.Student{
							{"BENST129", "BENJAMIN SANTOS", sql.NullString{}},
							{"YAPW9087", "YAP WEI LING", sql.NullString{}},
						},
					},
					{
						"A32",
						[]database.ClassVenueSchedule{},
						[]database.Student{
							{"PATELAR14", "ARJUN PATEL", sql.NullString{}},
							{"YAPX9087", "YAP XIN TING", sql.NullString{}},
						},
					},
				},
			},
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.file)
			a.Nil(err)

			data, err := ParseClassCreationFile(file)
			a.Nil(err)

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
