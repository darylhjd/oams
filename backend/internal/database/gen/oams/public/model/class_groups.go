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

type ClassGroups struct {
	ID        int64 `sql:"primary_key"`
	ClassID   int64
	Name      string
	ClassType ClassType
	CreatedAt time.Time
	UpdatedAt time.Time
}
