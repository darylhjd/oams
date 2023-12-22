package permissions

import "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"

type P int

const (
	BatchPost P = iota
	BatchPut

	UserRead

	ClassRead

	ClassManagerCreate
	ClassManagerRead

	ClassGroupRead

	ClassGroupSessionRead

	SessionEnrollmentRead
	SessionEnrollmentUpdate

	AttendanceTakingRead
)

type permissionMap map[P]struct{}

var rolePermissionMapping = map[model.UserRole]permissionMap{
	model.UserRole_User:        userRolePermissions,
	model.UserRole_SystemAdmin: systemAdminRolePermissions,
}

var userRolePermissions = permissionMap{
	UserRead:                {},
	ClassRead:               {},
	ClassGroupRead:          {},
	ClassGroupSessionRead:   {},
	SessionEnrollmentRead:   {},
	SessionEnrollmentUpdate: {},
	AttendanceTakingRead:    {},
}

var systemAdminRolePermissions = permissionMap{
	BatchPost: {},
	BatchPut:  {},

	UserRead: {},

	ClassRead: {},

	ClassManagerCreate: {},
	ClassManagerRead:   {},

	ClassGroupRead: {},

	ClassGroupSessionRead: {},

	SessionEnrollmentRead:   {},
	SessionEnrollmentUpdate: {},

	AttendanceTakingRead: {},
}
