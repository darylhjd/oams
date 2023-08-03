package common

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/xuri/excelize/v2"

	"github.com/darylhjd/oams/backend/internal/database"
)

const (
	namespace = "apiserver/common"
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

// ParseBatchFile parses a class creation file.
func ParseBatchFile(filename string, f io.Reader) (BatchData, error) {
	file, err := excelize.OpenReader(f)
	if err != nil {
		return BatchData{}, fmt.Errorf("%s - cannot open request file: %w", namespace, err)
	}
	defer func() {
		_ = file.Close()
	}()

	var creationData BatchData
	creationData.Filename = filename

	sheets := file.GetSheetList()
	if len(sheets) != expectedSheetCount {
		return creationData, fmt.Errorf("%s - invalid class creation file format", namespace)
	}

	rows, err := file.GetRows(sheets[expectedSheetCount-1])
	if err != nil {
		return creationData, fmt.Errorf("%s - cannot get data rows", namespace)
	}

	if err = parseClassMetaData(&creationData, rows); err != nil {
		return creationData, fmt.Errorf("%s - error while parsing class metadata: %w", namespace, err)
	}

	if err = parseClassGroups(&creationData, rows); err != nil {
		return creationData, fmt.Errorf("%s - error while parsing class groups: %w", namespace, err)
	}

	return creationData, creationData.IsValid()
}

// parseClassMetaData is a helper function to parse a class' metadata from a file.
func parseClassMetaData(creationData *BatchData, rows [][]string) error {
	// The first few rows in the sheet should be the course metadata.
	if len(rows) < expectedClassMetaDataRows {
		return errors.New("not enough rows for class metadata")
	}

	// Each metadata row has an expected number of filled columns.
	for i := 0; i < expectedClassMetaDataRows; i++ {
		if len(rows[i]) != expectedClassMetaDataRowLength {
			return fmt.Errorf("unexpected number of columns for class metadata row %d", i+1)
		}
	}

	// Parse class Year and Semester, as well as the time of creation of the class creation file.
	var (
		d string
		t string
	)
	if _, err := fmt.Sscanf(rows[yearSemesterDateRow][expectedClassMetaDataRowLength-1], yearSemesterDateFormat,
		&creationData.Course.Year, &creationData.Course.Semester, &d, &t); err != nil {
		return fmt.Errorf("could not parse class year and semester: %w", err)
	}

	loc, err := time.LoadLocation(timezoneLocation)
	if err != nil {
		return fmt.Errorf("invalid timezone name when parsing class creation file creation date: %w", err)
	}

	date, err := time.ParseInLocation(creationDateFormat, fmt.Sprintf("%s %s", d, t), loc)
	if err != nil {
		return fmt.Errorf("could not parse class creation file creation date: %w", err)
	}

	creationData.FileCreationDate = date

	// Parse course programme.
	creationData.Course.Programme = strings.TrimPrefix(rows[courseProgrammeRow][expectedClassMetaDataRowLength-1], courseProgrammePrefix)

	// Parse class course code.
	if _, err = fmt.Sscanf(rows[courseCodeRow][expectedClassMetaDataRowLength-1], courseCodeFormat,
		&creationData.Course.Code, &creationData.Course.Au); err != nil {
		return fmt.Errorf("could not parse course code and au count: %w", err)
	}

	// Parse class type.
	if _, err = fmt.Sscanf(rows[classTypeRow][expectedClassMetaDataRowLength-1], classTypeFormat,
		&creationData.classType); err != nil {
		return fmt.Errorf("could not parse class type: %w", err)
	}

	return nil
}

// parseClassGroups is a helper function to parse a class' groups.
func parseClassGroups(creationData *BatchData, rows [][]string) error {
	index := expectedClassMetaDataRows + 1            // Skip one blank row after metadata.
	for index+expectedClassGroupIDRows <= len(rows) { // For each class group.
		var group ClassGroupData
		group.ClassType = creationData.classType

		// Parse class group ID.
		if len(rows[index]) != expectedClassGroupMetaDataRowLength {
			return errors.New("unexpected number of columns for class group row")
		}

		if _, err := fmt.Sscanf(rows[index][expectedClassGroupMetaDataRowLength-1], classGroupFormat,
			&group.Name); err != nil {
			return fmt.Errorf("could not parse class group: %w", err)
		}

		// Parse class group sessions.
		index += expectedClassGroupIDRows
		for index+expectedClassGroupSessionRows <= len(rows) && len(rows[index]) != 0 { // For each session for a class group.
			var (
				session   SessionData
				dayOfWeek string
				from      string
				to        string
				weeks     string
			)

			// Parse session day-time.
			if _, err := fmt.Sscanf(rows[index][expectedClassGroupMetaDataRowLength-1], classGroupSessionDayTimeFormat,
				&dayOfWeek, &from, &to, &weeks); err != nil {
				return fmt.Errorf("could not parse class group day-time: %w", err)
			}

			// TODO: Store actual date with the time.
			startTime, err := time.Parse(classGroupSessionTimeFormat, from)
			if err != nil {
				return fmt.Errorf("could not parse class group session start time: %w", err)
			}

			session.StartTime = pgtype.Timestamp{
				Time:  startTime,
				Valid: true,
			}

			endTime, err := time.Parse(classGroupSessionTimeFormat, to)
			if err != nil {
				return fmt.Errorf("could not parse class group session end time: %w", err)
			}

			session.EndTime = pgtype.Timestamp{
				Time:  endTime,
				Valid: true,
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
			return errors.New("unexpected start of class group enrollment list")
		}

		index += 1
		for index+1 <= len(rows) && len(rows[index]) != 0 {
			if len(rows[index]) != expectedStudentEnrollmentRowLength {
				return errors.New("unexpected number of columns for student enrollment row")
			}

			group.Students = append(group.Students, database.UpsertUsersParams{
				ID:   rows[index][studentIdColumn],
				Name: rows[index][studentNameColumn],
				Role: database.UserRoleSTUDENT,
			})

			index += 1
		}

		creationData.ClassGroups = append(creationData.ClassGroups, group)
		index += 1 // Skip blank row after end of enrollment list. Will also work if it is the last list.
	}

	return nil
}
