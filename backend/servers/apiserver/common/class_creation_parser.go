package common

import (
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/darylhjd/oams/backend/internal/database"
)

const (
	namespace = "apiserver/common"
)

const (
	timezoneLocation    = "Asia/Singapore"
	yearSemesterDateRow = 0
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
	courseCodeFormat               = "Course: %s %dAU"
	classTypeFormat                = "Class Type: %s"
	classGroupFormat               = "Class Group: %s"
	classGroupSessionDayTimeFormat = "Day-Time: %s  %s To: %s Wk%s"
	classGroupSessionTimeFormat    = "1504"
	classGroupSessionVenuePrefix   = "Venue: "
	classGroupEnrollmentListIdent  = "No."
)

// ClassCreationData is a struct containing data for a class creation file.
type ClassCreationData struct {
	Year             int
	Semester         string
	FileCreationDate time.Time // Time of file generation by NTU systems.
	CourseCode       string
	ClassType        string
	ClassGroups      []ClassGroup
}

// ClassGroup is a struct containing data for a class group for a class.
type ClassGroup struct {
	Name     string
	Sessions []database.ClassVenueSchedule
	Students []database.Student
}

// ParseClassCreationFile parses a class create file.
func ParseClassCreationFile(f io.Reader) (*ClassCreationData, error) {
	file, err := excelize.OpenReader(f)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	sheets := file.GetSheetList()
	if len(sheets) != expectedSheetCount {
		return nil, fmt.Errorf("%s - invalid class creation file format", namespace)
	}

	rows, err := file.GetRows(sheets[expectedSheetCount-1])
	if err != nil {
		return nil, fmt.Errorf("%s - cannot get data rows", namespace)
	}

	var classData ClassCreationData
	if err := parseClassMetaData(&classData, rows); err != nil {
		return nil, err
	}

	if err := parseClassGroups(&classData, rows); err != nil {
		return nil, err
	}

	return &classData, nil
}

// parseClassMetaData is a helper function to parse a class' metadata from a file.
func parseClassMetaData(classData *ClassCreationData, rows [][]string) error {
	// The first few rows in the sheet should be the course metadata.
	if len(rows) < expectedClassMetaDataRows {
		return fmt.Errorf("%s - not enough rows for class metadata", namespace)
	}

	// Each metadata row has an expected number of filled columns.
	for i := 0; i < expectedClassMetaDataRows; i++ {
		if len(rows[i]) != expectedClassMetaDataRowLength {
			return fmt.Errorf("%s - unexpected number of columns for class metadata row %d", namespace, i+1)
		}
	}

	// Parse class Year and Semester, as well as the time of creation of the class creation file.
	var (
		d string
		t string
	)
	if _, err := fmt.Sscanf(rows[yearSemesterDateRow][expectedClassMetaDataRowLength-1], yearSemesterDateFormat,
		&classData.Year, &classData.Semester, &d, &t); err != nil {
		return fmt.Errorf("%s - could not parse class year and semester: %w", namespace, err)
	}

	loc, err := time.LoadLocation(timezoneLocation)
	if err != nil {
		return fmt.Errorf("%s - invalid timezone name when parsing class creation file creation date: %w", namespace, err)
	}

	date, err := time.ParseInLocation(creationDateFormat, fmt.Sprintf("%s %s", d, t), loc)
	if err != nil {
		return fmt.Errorf("%s - could not parse class creation file creation date: %w", namespace, err)
	}

	classData.FileCreationDate = date

	// Parse class course code.
	var au int
	if _, err = fmt.Sscanf(rows[courseCodeRow][expectedClassMetaDataRowLength-1], courseCodeFormat,
		&classData.CourseCode, &au); err != nil {
		return fmt.Errorf("%s - could not parse course code and au count: %w", namespace, err)
	}

	// Parse class type.
	if _, err = fmt.Sscanf(rows[classTypeRow][expectedClassMetaDataRowLength-1], classTypeFormat,
		&classData.ClassType); err != nil {
		return fmt.Errorf("%s - could not parse class type: %w", namespace, err)
	}

	return nil
}

// parseClassGroups is a helper function to parse a class' groups.
func parseClassGroups(classData *ClassCreationData, rows [][]string) error {
	index := expectedClassMetaDataRows + 1            // Skip one blank row after metadata.
	for index+expectedClassGroupIDRows <= len(rows) { // For each class group.
		var group ClassGroup

		// Parse class group ID.
		if len(rows[index]) != expectedClassGroupMetaDataRowLength {
			return fmt.Errorf("%s - unexpected number of columns for class group row", namespace)
		}

		if _, err := fmt.Sscanf(rows[index][expectedClassGroupMetaDataRowLength-1], classGroupFormat,
			&group.Name); err != nil {
			return fmt.Errorf("%s - could not parse class group: %w", namespace, err)
		}

		// Parse class group sessions.
		index += expectedClassGroupIDRows
		for index+expectedClassGroupSessionRows <= len(rows) && len(rows[index]) != 0 { // For each session for a class group.
			var (
				session   database.ClassVenueSchedule
				dayOfWeek string
				from      string
				to        string
				weeks     string
			)

			// Parse session day-time.
			if _, err := fmt.Sscanf(rows[index][expectedClassGroupMetaDataRowLength-1], classGroupSessionDayTimeFormat,
				&dayOfWeek, &from, &to, &weeks); err != nil {
				return fmt.Errorf("%s - could not parse class group day-time: %w", namespace, err)
			}

			// TODO: Store actual date with the time.
			_, err := time.Parse(classGroupSessionTimeFormat, from)
			if err != nil {
				return fmt.Errorf("%s - could not parse class group session start time: %w", namespace, err)
			}

			_, err = time.Parse(classGroupSessionTimeFormat, to)
			if err != nil {
				return fmt.Errorf("%s - could not parse class group session end time: %w", namespace, err)
			}

			// Parse session venue.
			session.Venue = strings.TrimPrefix(rows[index+1][expectedClassGroupMetaDataRowLength-1], classGroupSessionVenuePrefix)

			group.Sessions = append(group.Sessions, session)
			index += expectedClassGroupSessionRows
		}

		// Parse student enrollment list.
		index += 1 // Skip one blank row after metadata.
		// Perform sanity check.
		if index+1 > len(rows) ||
			len(rows[index]) != expectedClassGroupEnrollmentIdentRowLength ||
			rows[index][0] != classGroupEnrollmentListIdent {
			return fmt.Errorf("%s - unable to parse enrollment list for class group", namespace)
		}

		index += 1
		for index+1 <= len(rows) && len(rows[index]) != 0 {
			if len(rows[index]) != expectedStudentEnrollmentRowLength {
				return fmt.Errorf("%s - unexpected number of columns for student enrollment row", namespace)
			}

			group.Students = append(group.Students, database.Student{
				ID:    rows[index][studentIdColumn],
				Name:  rows[index][studentNameColumn],
				Email: sql.NullString{},
			})

			index += 1
		}

		classData.ClassGroups = append(classData.ClassGroups, group)

		index += 1 // Skip blank row after end of enrollment list. Will also work if it is the last list.
	}

	if len(classData.ClassGroups) == 0 {
		return fmt.Errorf("%s - class creation file does not specify any valid class groups", namespace)
	}

	return nil
}
