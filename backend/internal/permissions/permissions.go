package permissions

type Permission int

const (
	PermissionUserCreate Permission = iota
	PermissionUserRead
	PermissionUserUpdate
	PermissionUserDelete

	PermissionClassCreate
	PermissionClassRead
)

// UserRolePermissions holds the default permission model for a User.
var UserRolePermissions = map[Permission]struct{}{
	PermissionUserRead:  {},
	PermissionClassRead: {},
}
