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

type ClassGroupSessions struct {
	ID           int64 `sql:"primary_key"`
	ClassGroupID int64
	StartTime    time.Time
	EndTime      time.Time
	Venue        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
