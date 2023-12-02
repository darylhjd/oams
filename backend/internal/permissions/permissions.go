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

	ClassManagerCreate
	ClassManagerRead
	ClassManagerUpdate
	ClassManagerDelete

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

func GetPermissions(role model.UserRole) PermissionMap {
	return rolePermissionMapping[role]
}

type PermissionMap map[P]struct{}

var rolePermissionMapping = map[model.UserRole]PermissionMap{
	model.UserRole_User: userRolePermissions,
}

var userRolePermissions = PermissionMap{
	UserRead:              {},
	ClassRead:             {},
	ClassGroupRead:        {},
	ClassGroupSessionRead: {},
	SessionEnrollmentRead: {},
}
