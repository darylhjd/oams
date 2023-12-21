package permissions

import "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"

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

	AttendanceTakingRead
)

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
	AttendanceTakingRead:  {},
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

	ClassManagerCreate: {},
	ClassManagerRead:   {},
	ClassManagerUpdate: {},
	ClassManagerDelete: {},

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

	AttendanceTakingRead: {},
}
