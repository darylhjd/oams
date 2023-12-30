//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/darylhjd/oams/backend/internal/database/types"
	"time"
)

type ClassAttendanceRule struct {
	ID          int64             `sql:"primary_key" json:"id"`
	ClassID     int64             `json:"class_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	RuleType    RuleType          `json:"rule_type"`
	Rule        string            `json:"rule"`
	Environment types.Environment `json:"environment"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
