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

	// classType allows each class group to have access to class type information during processing.
	// The class type information is discovered only when processing a class' metadata, which only occurs once
	// and happens before processing for class groups can begin.
	classType database.ClassType
}

// ClassGroupData is a struct containing data for creating a class group and its associated sessions and students.
type ClassGroupData struct {
	database.UpsertClassGroupsParams
	Sessions []database.UpsertClassGroupSessionsParams `json:"sessions"`
	Students []database.UpsertUsersParams              `json:"students"`
}
