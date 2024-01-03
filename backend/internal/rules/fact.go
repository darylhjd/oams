package rules

import "time"

// Fact represents a session attendance fact. Each fact corresponds to one session. Additional useful information is also
// contained within to allow the user to generate custom rules.
type Fact struct {
	ClassID       int64     `alias:"class.id"`
	ClassCode     string    `alias:"class.code"`
	ClassYear     int32     `alias:"class.year"`
	ClassSemester string    `alias:"class.semester"`
	ClassType     string    `alias:"class_group.class_type"`
	StartTime     time.Time `alias:"class_group_session.start_time"`
	EndTime       time.Time `alias:"class_group_session.end_time"`
	Venue         string    `alias:"class_group_session.venue"`
	UserID        string    `alias:"session_enrollment.user_id"`
	UserName      string    `alias:"user.name"`
	UserEmail     string    `alias:"user.email"`
	Attended      bool      `alias:"session_enrollment.attended"`
}
