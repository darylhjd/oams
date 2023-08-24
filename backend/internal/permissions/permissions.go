package permissions

import "github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"

type Permission int

const (
	PermissionBatchPost Permission = iota
	PermissionBatchPut

	PermissionUserCreate
	PermissionUserRead
	PermissionUserUpdate
	PermissionUserDelete

	PermissionClassCreate
	PermissionClassRead
	PermissionClassUpdate
	PermissionClassDelete

	PermissionClassGroupCreate
	PermissionClassGroupRead
	PermissionClassGroupUpdate
	PermissionClassGroupDelete

	PermissionClassGroupSessionCreate
	PermissionClassGroupSessionRead
	PermissionClassGroupSessionUpdate
	PermissionClassGroupSessionDelete

	PermissionSessionEnrollmentCreate
	PermissionSessionEnrollmentRead
	PermissionSessionEnrollmentUpdate
	PermissionSessionEnrollmentDelete
)

// HasPermissions checks if a user with a role has all the given permissions.
func HasPermissions(role model.UserRole, permissions ...Permission) bool {
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

type permissionMap map[Permission]struct{}

var rolePermissionMapping = map[model.UserRole]permissionMap{
	model.UserRole_User:        userRolePermissions,
	model.UserRole_SystemAdmin: systemAdminRolePermissions,
}

var userRolePermissions = permissionMap{
	PermissionUserRead:              {},
	PermissionClassRead:             {},
	PermissionClassGroupRead:        {},
	PermissionClassGroupSessionRead: {},
	PermissionSessionEnrollmentRead: {},
}

var systemAdminRolePermissions = permissionMap{
	PermissionBatchPost: {},
	PermissionBatchPut:  {},

	PermissionUserCreate: {},
	PermissionUserRead:   {},
	PermissionUserUpdate: {},
	PermissionUserDelete: {},

	PermissionClassCreate: {},
	PermissionClassRead:   {},
	PermissionClassUpdate: {},
	PermissionClassDelete: {},

	PermissionClassGroupCreate: {},
	PermissionClassGroupRead:   {},
	PermissionClassGroupUpdate: {},
	PermissionClassGroupDelete: {},

	PermissionClassGroupSessionCreate: {},
	PermissionClassGroupSessionRead:   {},
	PermissionClassGroupSessionUpdate: {},
	PermissionClassGroupSessionDelete: {},

	PermissionSessionEnrollmentCreate: {},
	PermissionSessionEnrollmentRead:   {},
	PermissionSessionEnrollmentUpdate: {},
	PermissionSessionEnrollmentDelete: {},
}
