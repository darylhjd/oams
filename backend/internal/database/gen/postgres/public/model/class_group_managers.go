//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type ClassGroupManager struct {
	ID           int64        `sql:"primary_key" json:"id"`
	UserID       string       `json:"user_id"`
	ClassGroupID int64        `json:"class_group_id"`
	ManagingRole ManagingRole `json:"managing_role"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
