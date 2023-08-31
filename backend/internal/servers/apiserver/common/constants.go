package common

import (
	"log"
	"time"
	_ "time/tzdata"
)

const (
	timezoneLocation    = "Asia/Singapore"
	yearSemesterDateRow = 0
	courseProgrammeRow  = 1
	courseCodeRow       = 2
	classTypeRow        = 3
	studentNameColumn   = 1
	studentIdColumn     = 5

	expectedSheetCount                         = 1
	expectedClassMetaDataRows                  = 4
	expectedClassMetaDataRowLength             = 1
	expectedClassGroupIDRows                   = 1
	expectedClassGroupSessionRows              = 2
	expectedClassGroupMetaDataRowLength        = 1
	expectedClassGroupEnrollmentIdentRowLength = 3
	expectedStudentEnrollmentRowLength         = 6
)

const (
	yearSemesterDateFormat         = "Class Attendance List: %d, %s  Date: %s %s"
	creationDateFormat             = "02-Jan-2006 15:04"
	courseProgrammePrefix          = "Programme: "
	courseCodeFormat               = "Course: %s %dAU"
	classTypeFormat                = "Class Type: %s"
	classGroupFormat               = "Class Group: %s"
	classGroupSessionDayTimeFormat = "Day-Time: %s  %s To: %s Wk%s"
	classGroupSessionTimeFormat    = "1504"
	classGroupSessionVenuePrefix   = "Venue: "
	classGroupEnrollmentListIdent  = "No."
)

const (
	semester1 = "1"
	semester2 = "2"

	classGroupSessionWeekCommaSep             = ","
	classGroupSessionWeekHyphenSep            = "-"
	classGroupSessionWeekHyphenExpectedLength = 2

	semester1YearWeek   = 33
	semester2YearWeek   = 2
	recessWeekAfterWeek = 7
)

var location *time.Location

func init() {
	log.Println("init called in common package")
	var err error
	location, err = time.LoadLocation(timezoneLocation)
	if err != nil {
		log.Printf("Error: %s", err)
	}
}
