package permissions

import "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"

type P int

const (
	BatchPost P = iota
	BatchPut

	UserRead
	UserUpdate

	ClassRead

	ClassManagerCreate
	ClassManagerRead

	ClassGroupRead

	ClassGroupSessionRead

	SessionEnrollmentRead
	SessionEnrollmentUpdate

	AttendanceTakingRead
	AttendanceTakingUpdate
)

type permissionMap map[P]struct{}

var rolePermissionMapping = map[model.UserRole]permissionMap{
	model.UserRole_User:        userRolePermissions,
	model.UserRole_SystemAdmin: systemAdminRolePermissions,
}

var userRolePermissions = permissionMap{
	UserRead:                {},
	UserUpdate:              {},
	ClassRead:               {},
	ClassGroupRead:          {},
	ClassGroupSessionRead:   {},
	SessionEnrollmentRead:   {},
	SessionEnrollmentUpdate: {},
	AttendanceTakingRead:    {},
	AttendanceTakingUpdate:  {},
}

var systemAdminRolePermissions = permissionMap{
	BatchPost: {},
	BatchPut:  {},

	UserRead:   {},
	UserUpdate: {},

	ClassRead: {},

	ClassManagerCreate: {},
	ClassManagerRead:   {},

	ClassGroupRead: {},

	ClassGroupSessionRead: {},

	SessionEnrollmentRead:   {},
	SessionEnrollmentUpdate: {},

	AttendanceTakingRead:   {},
	AttendanceTakingUpdate: {},
}
