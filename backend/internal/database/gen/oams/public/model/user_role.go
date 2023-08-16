//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type UserRole string

const (
	UserRole_Student           UserRole = "STUDENT"
	UserRole_CourseCoordinator UserRole = "COURSE_COORDINATOR"
	UserRole_Admin             UserRole = "ADMIN"
)

func (e *UserRole) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("jet: Invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "STUDENT":
		*e = UserRole_Student
	case "COURSE_COORDINATOR":
		*e = UserRole_CourseCoordinator
	case "ADMIN":
		*e = UserRole_Admin
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for UserRole enum")
	}

	return nil
}

func (e UserRole) String() string {
	return string(e)
}
