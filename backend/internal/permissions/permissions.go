package permissions

import "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"

type P int

const (
	BatchPost P = iota
	BatchPut

	UserCreate
	UserRead
	UserUpdate
	UserDelete

	ClassCreate
	ClassRead
	ClassUpdate
	ClassDelete

	ClassGroupCreate
	ClassGroupRead
	ClassGroupUpdate
	ClassGroupDelete

	ClassGroupSessionCreate
	ClassGroupSessionRead
	ClassGroupSessionUpdate
	ClassGroupSessionDelete

	SessionEnrollmentCreate
	SessionEnrollmentRead
	SessionEnrollmentUpdate
	SessionEnrollmentDelete
)

// HasPermissions checks if a user with a role has all the given permissions.
func HasPermissions(role model.UserRole, permissions ...P) bool {
	permModel, ok := rolePermissionMapping[role]
	if !ok {
		return false
	}

	for _, perm := range permissions {
		if _, ok = permModel[perm]; !ok {
			return false
		}
	}

	return true
}

type permissionMap map[P]struct{}

var rolePermissionMapping = map[model.UserRole]permissionMap{
	model.UserRole_User:        userRolePermissions,
	model.UserRole_SystemAdmin: systemAdminRolePermissions,
}

var userRolePermissions = permissionMap{
	UserRead:              {},
	ClassRead:             {},
	ClassGroupRead:        {},
	ClassGroupSessionRead: {},
	SessionEnrollmentRead: {},
}

var systemAdminRolePermissions = permissionMap{
	BatchPost: {},
	BatchPut:  {},

	UserCreate: {},
	UserRead:   {},
	UserUpdate: {},
	UserDelete: {},

	ClassCreate: {},
	ClassRead:   {},
	ClassUpdate: {},
	ClassDelete: {},

	ClassGroupCreate: {},
	ClassGroupRead:   {},
	ClassGroupUpdate: {},
	ClassGroupDelete: {},

	ClassGroupSessionCreate: {},
	ClassGroupSessionRead:   {},
	ClassGroupSessionUpdate: {},
	ClassGroupSessionDelete: {},

	SessionEnrollmentCreate: {},
	SessionEnrollmentRead:   {},
	SessionEnrollmentUpdate: {},
	SessionEnrollmentDelete: {},
}
