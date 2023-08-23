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
)

// UserRolePermissions holds the default permission model for a User.
var UserRolePermissions = map[Permission]struct{}{
	PermissionUserRead: {},
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
}

// HasPermission checks if a user with a role has all the given permissions.
func HasPermission(role model.UserRole, permissions ...Permission) bool {
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
