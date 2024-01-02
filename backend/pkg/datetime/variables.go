package datetime

import "time"

const (
	timezone = "Asia/Singapore"
)

var (
	weekdays = map[string]time.Weekday{}

	// Location contains the time.Location object for Asia/Singapore
	Location *time.Location
)

func init() {
	for d := time.Sunday; d <= time.Saturday; d++ {
		name := d.String()
		weekdays[name] = d
		weekdays[name[:3]] = d
	}

	Location, _ = time.LoadLocation(timezone)
}
