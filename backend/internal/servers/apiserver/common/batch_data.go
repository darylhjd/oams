package common

import (
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
)

// BatchData is a struct containing data for processing a batch of entities.
type BatchData struct {
	// Filename and FileCreationDate is specific to data generated from a file upload.
	Filename         string    `json:"filename"`
	FileCreationDate time.Time `json:"file_creation_date"`

	Class       database.UpsertClassesParams `json:"class"`
	ClassGroups []ClassGroupData             `json:"class_groups"`

	// Helper fields for filling data during processing.
	classType database.ClassType
}

// ClassGroupData is a struct containing data for creating a class group and its associated sessions.
type ClassGroupData struct {
	database.UpsertClassGroupsParams
	Sessions []SessionData                `json:"sessions"`
	Students []database.UpsertUsersParams `json:"students"`
}

// SessionData is a struct containing information on a session.
type SessionData struct {
	database.UpsertClassGroupSessionsParams

	// Helper field for filling session data.
	Course *BatchData `json:"-"`
}
