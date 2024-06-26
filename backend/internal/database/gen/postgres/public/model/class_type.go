//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type ClassType string

const (
	ClassType_Lec ClassType = "LEC"
	ClassType_Tut ClassType = "TUT"
	ClassType_Lab ClassType = "LAB"
)

func (e *ClassType) Scan(value interface{}) error {
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
	case "LEC":
		*e = ClassType_Lec
	case "TUT":
		*e = ClassType_Tut
	case "LAB":
		*e = ClassType_Lab
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for ClassType enum")
	}

	return nil
}

func (e ClassType) String() string {
	return string(e)
}
