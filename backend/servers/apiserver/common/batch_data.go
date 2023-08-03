package common

import (
	"errors"
	"fmt"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
)

// BatchData is a struct containing data for creating a class. It is used by both file and JSON requests
// for class creation. Prior to processing the data here, the caller should call IsValid on the data at least once.
type BatchData struct {
	// Filename and FileCreationDate is specific to creation data generated by file upload.
	Filename         string    `json:"filename"`
	FileCreationDate time.Time `json:"file_creation_date"`

	Course      database.UpsertClassesParams `json:"course"`
	ClassGroups []ClassGroupData             `json:"class_groups"`

	// These variables are helper fields for generating filled data.
	classType database.ClassType
}

// ClassGroupData is a struct containing data for creating a class group and its associated sessions.
type ClassGroupData struct {
	database.UpsertClassGroupsParams
	Sessions []SessionData                `json:"sessions"`
	Students []database.UpsertUsersParams `json:"students"`
}

type SessionData struct {
	Course *BatchData `json:"-"`
	database.UpsertClassGroupSessionsParams
}

// IsValid is a helper function to check if a BatchData is valid and returns an error if it is not.
func (c *BatchData) IsValid() error {
	if len(c.ClassGroups) == 0 {
		return errors.New("creation data has no valid class groups")
	}

	for _, classGroup := range c.ClassGroups {
		if len(classGroup.Sessions) == 0 {
			return fmt.Errorf("class group %s has no sessions", classGroup.Name)
		}

		if len(classGroup.Students) == 0 {
			return fmt.Errorf("class group %s has no enrollments", classGroup.Name)
		}
	}

	return nil
}
