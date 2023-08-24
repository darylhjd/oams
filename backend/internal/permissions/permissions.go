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
	var permModel map[Permission]struct{}
	switch role {
	case model.UserRole_User:
		permModel = UserRolePermissions
	case model.UserRole_SystemAdmin:
		permModel = SystemAdminRolePermissions
	default:
		return false
	}

	for _, perm := range permissions {
		if _, ok := permModel[perm]; !ok {
			return false
		}
	}

	return true
}

// UserRolePermissions holds the default permission model for a User.
var UserRolePermissions = map[Permission]struct{}{
	PermissionUserRead:              {},
	PermissionClassRead:             {},
	PermissionClassGroupRead:        {},
	PermissionClassGroupSessionRead: {},
	PermissionSessionEnrollmentRead: {},
}

var SystemAdminRolePermissions = map[Permission]struct{}{
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
