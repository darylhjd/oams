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

type Class struct {
	ID        int64     `sql:"primary_key" json:"id"`
	Code      string    `json:"code"`
	Year      int32     `json:"year"`
	Semester  string    `json:"semester"`
	Programme string    `json:"programme"`
	Au        int16     `json:"au"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
