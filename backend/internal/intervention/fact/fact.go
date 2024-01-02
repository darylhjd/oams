// Package fact cannot import database model package.
package fact

import "time"

// F represents an enrollment fact. Each fact corresponds to one session. Additional useful information is also
// contained within to allow the user to generate custom rules.
type F struct {
	ClassID   int64     `alias:"class.id"`
	StartTime time.Time `alias:"class_group_session.start_time"`
	EndTime   time.Time `alias:"class_group_session.end_time"`
	Venue     string    `alias:"class_group_session.venue"`
	UserID    string    `alias:"session_enrollment.user_id"`
	Attended  bool      `alias:"session_enrollment.attended"`
}
