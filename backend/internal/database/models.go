// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package database

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type ClassType string

const (
	ClassTypeLEC ClassType = "LEC"
	ClassTypeTUT ClassType = "TUT"
	ClassTypeLAB ClassType = "LAB"
)

func (e *ClassType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ClassType(s)
	case string:
		*e = ClassType(s)
	default:
		return fmt.Errorf("unsupported scan type for ClassType: %T", src)
	}
	return nil
}

type NullClassType struct {
	ClassType ClassType `json:"class_type"`
	Valid     bool      `json:"valid"` // Valid is true if ClassType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullClassType) Scan(value interface{}) error {
	if value == nil {
		ns.ClassType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ClassType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullClassType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ClassType), nil
}

type ClassGroup struct {
	ID        int64            `json:"id"`
	CourseID  int64            `json:"course_id"`
	Name      string           `json:"name"`
	ClassType ClassType        `json:"class_type"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type ClassGroupSession struct {
	ID           int64            `json:"id"`
	ClassGroupID int64            `json:"class_group_id"`
	StartTime    pgtype.Timestamp `json:"start_time"`
	EndTime      pgtype.Timestamp `json:"end_time"`
	Venue        string           `json:"venue"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
	UpdatedAt    pgtype.Timestamp `json:"updated_at"`
}

type Course struct {
	ID        int64            `json:"id"`
	Code      string           `json:"code"`
	Year      int32            `json:"year"`
	Semester  string           `json:"semester"`
	Programme string           `json:"programme"`
	Au        int16            `json:"au"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type SessionEnrollment struct {
	SessionID int64  `json:"session_id"`
	StudentID string `json:"student_id"`
}

type Student struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Email     pgtype.Text      `json:"email"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}
