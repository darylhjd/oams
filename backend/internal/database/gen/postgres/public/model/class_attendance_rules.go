//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/darylhjd/oams/backend/internal/rules"
	"time"
)

type ClassAttendanceRule struct {
	ID          int64             `sql:"primary_key" json:"id"`
	ClassID     int64             `json:"class_id"`
	CreatorID   string            `json:"creator_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Rule        string            `json:"rule"`
	Environment rules.Environment `json:"environment"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
