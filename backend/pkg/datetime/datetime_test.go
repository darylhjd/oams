package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeekStart(t *testing.T) {
	tts := []struct {
		name         string
		year         int
		weekNumber   int
		expectedDate string // Expected date in "2006-01-02" format.
	}{
		{"1st week of 2023", 2023, 1, "2023-01-02"},
		{"2nd week of 2023", 2023, 2, "2023-01-09"},
		{"3rd week of 2023", 2023, 3, "2023-01-16"},
		{"52nd week of 2023", 2023, 52, "2023-12-25"},
		{"last week of 2022", 2023, 0, "2022-12-26"},
		{"first week of 2024", 2023, 53, "2024-01-01"},
		{"first week of 2024", 2024, 1, "2024-01-01"},
		{"first week of 2023", 2022, 53, "2023-01-02"},
		{"first week of 2025", 2024, 54, "2025-01-06"},
		{"first week of 2001", 2000, 53, "2001-01-01"},
		{"first week of 2101", 2100, 53, "2101-01-03"},
		{"first week of 2002", 2001, 54, "2002-01-07"},
		{"2nd last week of 2000", 2001, -1, "2000-12-18"},
		{"3rd last week of 2000", 2001, -2, "2000-12-11"},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			expectedDate, err := time.ParseInLocation("2006-01-02", tt.expectedDate, time.Local)
			a.Nil(err)
			a.Equal(expectedDate, WeekStart(tt.year, tt.weekNumber, time.Local))
		})
	}
}

func TestParseWeekday(t *testing.T) {
	tts := []struct {
		name        string
		inputs      []string
		wantWeekday time.Weekday
		wantErr     bool
	}{
		{"Monday", []string{"Monday", "Mon"}, time.Monday, false},
		{"Tuesday", []string{"Tuesday", "Tue"}, time.Tuesday, false},
		{"Wednesday", []string{"Wednesday", "Wed"}, time.Wednesday, false},
		{"Thursday", []string{"Thursday", "Thu"}, time.Thursday, false},
		{"Friday", []string{"Friday", "Fri"}, time.Friday, false},
		{"Saturday", []string{"Saturday", "Sat"}, time.Saturday, false},
		{"Sunday", []string{"Sunday", "Sun"}, time.Sunday, false},
		{"invalid", []string{"invalid", "ii", "january"}, time.Sunday, true},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			for _, input := range tt.inputs {
				output, err := ParseWeekday(input)
				if tt.wantErr {
					a.Error(err)
				} else {
					a.Equal(tt.wantWeekday, output)
				}
			}
		})
	}
}
