//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/v2/postgres"

var UserRole = &struct {
	User        postgres.StringExpression
	SystemAdmin postgres.StringExpression
}{
	User:        postgres.NewEnumValue("USER"),
	SystemAdmin: postgres.NewEnumValue("SYSTEM_ADMIN"),
}
