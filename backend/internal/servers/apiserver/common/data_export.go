package common

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/darylhjd/oams/backend/internal/database"
)

// GenerateDataExportArchive generates a full CSV snapshot of the database within a zip file that
// is written to a http.ResponseWriter.
func GenerateDataExportArchive(w http.ResponseWriter, r *http.Request, txDb *database.DB) error {
	archive := zip.NewWriter(w)

	if err := generateUsersCSV(r, txDb, archive); err != nil {
		return fmt.Errorf("error generating users csv: %w", err)
	}

	if err := generateClassesCSV(r, txDb, archive); err != nil {
		return fmt.Errorf("error generating classes csv: %w", err)
	}

	if err := generateClassAttendanceRulesCSV(r, txDb, archive); err != nil {
		return fmt.Errorf("error generating class attendance rules csv: %w", err)
	}

	if err := generateClassGroupsCSV(r, txDb, archive); err != nil {
		return fmt.Errorf("error generating class groups csv: %w", err)
	}

	if err := generateClassGroupManagersCSV(r, txDb, archive); err != nil {
		return fmt.Errorf("error generating class group managers csv: %w", err)
	}

	if err := generateClassGroupSessionsCSV(r, txDb, archive); err != nil {
		return fmt.Errorf("error generating class group sessions csv: %w", err)
	}

	if err := generateSessionEnrollmentsCSV(r, txDb, archive); err != nil {
		return fmt.Errorf("error generating session enrollments csv: %w", err)
	}

	return archive.Close()
}

func generateUsersCSV(r *http.Request, txDb *database.DB, archive *zip.Writer) error {
	users, err := txDb.ListUsers(r.Context(), database.ListQueryParams{})
	if err != nil {
		return err
	}

	file, err := archive.Create("users.csv")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	records := append(
		make([][]string, 0, len(users)+1),
		[]string{"id", "name", "email", "role", "created_at", "updated_at"},
	)
	for _, user := range users {
		records = append(records, []string{
			user.ID,
			user.Name,
			user.Email,
			string(user.Role),
			user.CreatedAt.String(),
			user.UpdatedAt.String(),
		})
	}

	return writer.WriteAll(records)
}

func generateClassesCSV(r *http.Request, txDb *database.DB, archive *zip.Writer) error {
	classes, err := txDb.ListClasses(r.Context(), database.ListQueryParams{})
	if err != nil {
		return err
	}

	file, err := archive.Create("classes.csv")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	records := append(
		make([][]string, 0, len(classes)+1),
		[]string{"id", "code", "year", "semester", "programme", "au", "created_at", "updated_at"},
	)
	for _, class := range classes {
		records = append(records, []string{
			strconv.FormatInt(class.ID, 10),
			class.Code,
			strconv.FormatInt(int64(class.Year), 10),
			class.Semester,
			class.Programme,
			strconv.FormatInt(int64(class.Au), 10),
			class.CreatedAt.String(),
			class.UpdatedAt.String(),
		})
	}

	return writer.WriteAll(records)
}

func generateClassAttendanceRulesCSV(r *http.Request, txDb *database.DB, archive *zip.Writer) error {
	rules, err := txDb.ListClassAttendanceRules(r.Context(), database.ListQueryParams{})
	if err != nil {
		return err
	}

	file, err := archive.Create("class_attendance_rules.csv")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	records := append(
		make([][]string, 0, len(rules)+1),
		[]string{"id", "class_id", "creator_id", "title", "description", "rule", "environment", "created_at", "updated_at"},
	)
	for _, rule := range rules {
		environmentBytes, err := json.Marshal(rule.Environment)
		if err != nil {
			return err
		}

		records = append(records, []string{
			strconv.FormatInt(rule.ID, 10),
			strconv.FormatInt(rule.ClassID, 10),
			rule.CreatorID,
			rule.Title,
			rule.Description,
			rule.Rule,
			string(environmentBytes),
			rule.CreatedAt.String(),
			rule.UpdatedAt.String(),
		})
	}

	return writer.WriteAll(records)
}

func generateClassGroupsCSV(r *http.Request, txDb *database.DB, archive *zip.Writer) error {
	groups, err := txDb.ListClassGroups(r.Context(), database.ListQueryParams{})
	if err != nil {
		return err
	}

	file, err := archive.Create("class_groups.csv")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	records := append(
		make([][]string, 0, len(groups)+1),
		[]string{"id", "class_id", "name", "class_type", "created_at", "updated_at"},
	)
	for _, group := range groups {
		records = append(records, []string{
			strconv.FormatInt(group.ID, 10),
			strconv.FormatInt(group.ClassID, 10),
			group.Name,
			string(group.ClassType),
			group.CreatedAt.String(),
			group.UpdatedAt.String(),
		})
	}

	return writer.WriteAll(records)
}

func generateClassGroupManagersCSV(r *http.Request, txDb *database.DB, archive *zip.Writer) error {
	managers, err := txDb.ListClassGroupManagers(r.Context(), database.ListQueryParams{})
	if err != nil {
		return err
	}

	file, err := archive.Create("class_group_managers.csv")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	records := append(
		make([][]string, 0, len(managers)+1),
		[]string{"id", "user_id", "class_group_id", "managing_role", "created_at", "updated_at"},
	)
	for _, manager := range managers {
		records = append(records, []string{
			strconv.FormatInt(manager.ID, 10),
			manager.UserID,
			strconv.FormatInt(manager.ClassGroupID, 10),
			string(manager.ManagingRole),
			manager.CreatedAt.String(),
			manager.UpdatedAt.String(),
		})
	}

	return writer.WriteAll(records)
}

func generateClassGroupSessionsCSV(r *http.Request, txDb *database.DB, archive *zip.Writer) error {
	classGroups, err := txDb.ListClassGroupSessions(r.Context(), database.ListQueryParams{})
	if err != nil {
		return err
	}

	file, err := archive.Create("class_group_sessions.csv")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	records := append(
		make([][]string, 0, len(classGroups)+1),
		[]string{"id", "class_group_id", "start_time", "end_time", "venue", "created_at", "updated_at"},
	)
	for _, group := range classGroups {
		records = append(records, []string{
			strconv.FormatInt(group.ID, 10),
			strconv.FormatInt(group.ClassGroupID, 10),
			group.StartTime.String(),
			group.EndTime.String(),
			group.Venue,
			group.CreatedAt.String(),
			group.UpdatedAt.String(),
		})
	}

	return writer.WriteAll(records)
}

func generateSessionEnrollmentsCSV(r *http.Request, txDb *database.DB, archive *zip.Writer) error {
	enrollments, err := txDb.ListSessionEnrollments(r.Context(), database.ListQueryParams{})
	if err != nil {
		return err
	}

	file, err := archive.Create("session_enrollments.csv")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	records := append(
		make([][]string, 0, len(enrollments)+1),
		[]string{"id", "session_id", "user_id", "attended", "created_at", "updated_at"},
	)
	for _, enrollment := range enrollments {
		records = append(records, []string{
			strconv.FormatInt(enrollment.ID, 10),
			strconv.FormatInt(enrollment.SessionID, 10),
			enrollment.UserID,
			strconv.FormatBool(enrollment.Attended),
			enrollment.CreatedAt.String(),
			enrollment.UpdatedAt.String(),
		})
	}

	return writer.WriteAll(records)
}
