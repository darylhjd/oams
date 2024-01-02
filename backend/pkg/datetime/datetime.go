package datetime

import (
	"fmt"
	"time"
)

// WeekStart returns a time.Time for the given week's Monday in a given year.
// If location is nil, WeekStart will panic.
func WeekStart(year, week int, location *time.Location) time.Time {
	// Start from the middle of the year.
	t := time.Date(year, 7, 1, 0, 0, 0, 0, location)

	// Roll back to Monday.
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks.
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

// ParseWeekday returns a time.Weekday corresponding to the given name.
// Both long names such as "Monday", "Friday" and short names such as "Mon", "Fri" are recognized.
func ParseWeekday(s string) (time.Weekday, error) {
	if d, ok := weekdays[s]; ok {
		return d, nil
	}

	return time.Sunday, fmt.Errorf("invalid weekday %q", s)
}
