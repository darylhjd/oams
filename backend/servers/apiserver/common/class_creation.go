package common

import (
	"fmt"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
)

// ClassCreationData is a struct containing data for creating a class. It is used by both file and JSON requests
// for class creation. Prior to processing the data here, the caller should call Validate on the data at least once.
type ClassCreationData struct {
	// Filename and FileCreationDate is specific to creation data generated by file upload.
	Filename         string    `json:"filename"`
	FileCreationDate time.Time `json:"file_creation_date"`

	Course      database.CreateCoursesParams `json:"course"`
	ClassGroups []ClassGroupData             `json:"class_groups"`

	// Error provides an explanation for any failures while processing the creation data.
	Error string `json:"error"`

	// These variables are helper fields for generating filled data.
	ClassType database.ClassType `json:"-"`
}

// Validate is a helper function to check if a ClassCreationData is valid and sets the error message if necessary.
// Returns true if dat is valid and false otherwise.
func (c *ClassCreationData) Validate() bool {
	if c.Error != "" {
		return false
	}

	if len(c.ClassGroups) == 0 {
		c.Error = fmt.Sprintf("%s - creation data has no valid class groups", namespace)
		return false
	}

	for _, classGroup := range c.ClassGroups {
		if len(classGroup.Sessions) == 0 {
			c.Error = fmt.Sprintf("%s - class group %s has no sessions", namespace, classGroup.Name)
			return false
		}

		if len(classGroup.Students) == 0 {
			c.Error = fmt.Sprintf("%s - class group %s has no enrollments", namespace, classGroup.Name)
			return false
		}
	}

	return true
}

// ClassGroupData is a struct containing data for creating a class group and its associated sessions.
type ClassGroupData struct {
	database.CreateClassGroupsParams
	Sessions []database.CreateClassGroupSessionsParams `json:"sessions"`
	Students []database.UpsertStudentsParams           `json:"students"`
}
