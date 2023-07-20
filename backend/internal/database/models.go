// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ClassGroup struct {
	ID         int64
	CourseCode string
	Year       int32
	Semester   string
	GroupName  string
}

type ClassVenueSchedule struct {
	ClassGroupID int64
	Datetime     pgtype.Timestamp
	Venue        string
}

type Course struct {
	Code      string
	Programme string
	Au        int16
}

type Student struct {
	ID    string
	Name  string
	Email pgtype.Text
}
