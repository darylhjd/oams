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

type SessionEnrollment struct {
	ID        int64     `sql:"primary_key" json:"id"`
	SessionID int64     `json:"session_id"`
	UserID    string    `json:"user_id"`
	Attended  bool      `json:"attended"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
