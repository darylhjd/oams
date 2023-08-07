package common

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/darylhjd/oams/backend/pkg/datetime"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/xuri/excelize/v2"

	"github.com/darylhjd/oams/backend/internal/database"
)

// ParseBatchFile parses a class creation file.
func ParseBatchFile(filename string, f io.Reader) (BatchData, error) {
	file, err := excelize.OpenReader(f)
	if err != nil {
		return BatchData{}, fmt.Errorf("cannot open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var creationData BatchData
	creationData.Filename = filename

	sheets := file.GetSheetList()
	if len(sheets) != expectedSheetCount {
		return creationData, errors.New("invalid class creation file format")
	}

	rows, err := file.GetRows(sheets[expectedSheetCount-1])
	if err != nil {
		return creationData, errors.New("cannot get data rows")
	}

	if err = parseClassMetaData(&creationData, rows); err != nil {
		return creationData, fmt.Errorf("error while parsing class metadata: %w", err)
	}

	if err = parseClassGroups(&creationData, rows); err != nil {
		return creationData, fmt.Errorf("error while parsing class groups: %w", err)
	}

	return creationData, nil
}

// parseClassMetaData is a helper function to parse a class' metadata from a file.
func parseClassMetaData(batchData *BatchData, rows [][]string) error {
	// The first few rows in the sheet should be the class metadata.
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
	{
		var d, t string
		if _, err := fmt.Sscanf(rows[yearSemesterDateRow][expectedClassMetaDataRowLength-1], yearSemesterDateFormat,
			&batchData.Class.Year, &batchData.Class.Semester, &d, &t); err != nil {
			return fmt.Errorf("could not parse class year and semester: %w", err)
		}

		date, err := time.ParseInLocation(creationDateFormat, fmt.Sprintf("%s %s", d, t), location)
		if err != nil {
			return fmt.Errorf("could not parse class creation file creation date: %w", err)
		}

		batchData.FileCreationDate = date
	}

	// Parse course programme.
	batchData.Class.Programme = strings.TrimPrefix(rows[courseProgrammeRow][expectedClassMetaDataRowLength-1], courseProgrammePrefix)

	// Parse class course code.
	if _, err := fmt.Sscanf(rows[courseCodeRow][expectedClassMetaDataRowLength-1], courseCodeFormat,
		&batchData.Class.Code, &batchData.Class.Au); err != nil {
		return fmt.Errorf("could not parse course code and au count: %w", err)
	}

	// Parse class type.
	if _, err := fmt.Sscanf(rows[classTypeRow][expectedClassMetaDataRowLength-1], classTypeFormat,
		&batchData.classType); err != nil {
		return fmt.Errorf("could not parse class type: %w", err)
	}

	return nil
}

// parseClassGroups is a helper function to parse a class' groups.
func parseClassGroups(batchData *BatchData, rows [][]string) error {
	index := expectedClassMetaDataRows + 1            // Skip blank row after metadata.
	for index+expectedClassGroupIDRows <= len(rows) { // For each class group.
		var group ClassGroupData
		group.ClassType = batchData.classType

		// Parse class group name.
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

			// Parse session venue.
			venue := strings.TrimPrefix(rows[index+1][expectedClassGroupMetaDataRowLength-1], classGroupSessionVenuePrefix)

			sessions, err := parseClassGroupSessions(batchData, dayOfWeek, from, to, weeks, venue)
			if err != nil {
				return fmt.Errorf("could not parse class group sessions: %w", err)
			}

			group.Sessions = append(group.Sessions, sessions...)
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

		index += 1 // Skip enrollment list column row.
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

		batchData.ClassGroups = append(batchData.ClassGroups, group)
		index += 1 // Skip blank row after end of enrollment list. Will also work if it is the last list.
	}

	return nil
}

// parseClassGroupSessions is a helper function to create the appropriate sessions for a given class group session.
func parseClassGroupSessions(batchData *BatchData, dayOfWeek, from, to, weeksStr, venue string) ([]database.UpsertClassGroupSessionsParams, error) {
	var firstSessionStartDateTime, firstSessionEndDateTime time.Time
	{
		var year, week int
		switch batchData.Class.Semester {
		case semester1:
			year = int(batchData.Class.Year)
			week = semester1YearWeek
		case semester2:
			year = int(batchData.Class.Year + 1)
			week = semester2YearWeek
		default:
			return nil, errors.New("cannot guess semester start week due to unknown semester value")
		}

		day, err := datetime.ParseWeekday(dayOfWeek)
		if err != nil {
			return nil, err
		}

		startTime, err := time.Parse(classGroupSessionTimeFormat, from)
		if err != nil {
			return nil, fmt.Errorf("could not parse session start time: %w", err)
		}

		endTime, err := time.Parse(classGroupSessionTimeFormat, to)
		if err != nil {
			return nil, fmt.Errorf("could not parse session end time: %w", err)
		}

		startHour, startMinute, _ := startTime.Clock()
		endHour, endMinute, _ := endTime.Clock()

		firstSessionDate := datetime.WeekStart(year, week, location).AddDate(0, 0, int(day)-1)

		firstSessionStartDateTime = firstSessionDate.
			Add(time.Hour*time.Duration(startHour) + time.Minute*time.Duration(startMinute))
		firstSessionEndDateTime = firstSessionDate.
			Add(time.Hour*time.Duration(endHour) + time.Minute*time.Duration(endMinute))
	}

	// Parse the week numbers. There are 2 cases:
	// 1. If separated by hyphen (e.g. 2-13), then every week including the start and end weeks included.
	// 2. If separated by commas (e.g. 2,4,6,8), then each individual week included.
	var weeks []int
	switch {
	case strings.Contains(weeksStr, classGroupSessionWeekHyphenSep):
		startEnd := strings.Split(weeksStr, classGroupSessionWeekHyphenSep)
		if len(startEnd) != classGroupSessionWeekHyphenExpectedLength {
			return nil, errors.New("unexpected week formatting with hyphen separator")
		}

		startWeek, err := strconv.Atoi(startEnd[0])
		if err != nil {
			return nil, errors.New("start week number is not actually a number")
		}

		endWeek, err := strconv.Atoi(startEnd[classGroupSessionWeekHyphenExpectedLength-1])
		if err != nil {
			return nil, errors.New("end week number is not actually a number")
		}

		for i := startWeek; i <= endWeek; i++ {
			weeks = append(weeks, i)
		}
	case strings.Contains(weeksStr, classGroupSessionWeekCommaSep):
		for _, w := range strings.Split(weeksStr, classGroupSessionWeekCommaSep) {
			wInt, err := strconv.Atoi(w)
			if err != nil {
				return nil, errors.New("week number is not actually a number")
			}

			weeks = append(weeks, wInt)
		}
	default:
		return nil, errors.New("unexpected week formatting")
	}

	// For calculating session dates, add offset of 1 week after recess week.
	for idx := range weeks {
		if weeks[idx] > recessWeekAfterWeek {
			weeks[idx] += 1
		}
	}

	// Create all sessions.
	sessions := make([]database.UpsertClassGroupSessionsParams, 0, len(weeks))
	for _, week := range weeks {
		daysToAdd := 7 * (week - 1) // Since week count starts from 1.
		sessions = append(sessions, database.UpsertClassGroupSessionsParams{
			StartTime: pgtype.Timestamptz{
				Time:  firstSessionStartDateTime.AddDate(0, 0, daysToAdd),
				Valid: true,
			},
			EndTime: pgtype.Timestamptz{
				Time:  firstSessionEndDateTime.AddDate(0, 0, daysToAdd),
				Valid: true,
			},
			Venue: venue,
		})
	}

	return sessions, nil
}
