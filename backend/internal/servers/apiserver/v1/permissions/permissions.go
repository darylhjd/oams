package permissions

import "github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"

type P int

const (
	SignaturePut P = iota

	BatchPost
	BatchPut

	UserRead
	UserUpdate

	ClassRead

	ClassAttendanceRulesRead

	ClassGroupManagerPost
	ClassGroupManagerPut
	ClassGroupManagerRead

	ClassGroupRead

	ClassGroupSessionRead

	SessionEnrollmentRead

	AttendanceTakingRead
	AttendanceTakingUpdate
)

type permissionMap map[P]struct{}

var rolePermissionMapping = map[model.UserRole]permissionMap{
	model.UserRole_User:        userRolePermissions,
	model.UserRole_SystemAdmin: systemAdminRolePermissions,
}

var userRolePermissions = permissionMap{
	SignaturePut: {},

	UserRead:   {},
	UserUpdate: {},

	ClassRead: {},

	ClassAttendanceRulesRead: {},

	ClassGroupRead: {},

	ClassGroupSessionRead: {},

	SessionEnrollmentRead: {},

	AttendanceTakingRead:   {},
	AttendanceTakingUpdate: {},
}

var systemAdminRolePermissions = permissionMap{
	SignaturePut: {},

	BatchPost: {},
	BatchPut:  {},

	UserRead:   {},
	UserUpdate: {},

	ClassRead: {},

	ClassAttendanceRulesRead: {},

	ClassGroupManagerPost: {},
	ClassGroupManagerPut:  {},
	ClassGroupManagerRead: {},

	ClassGroupRead: {},

	ClassGroupSessionRead: {},

	SessionEnrollmentRead: {},

	AttendanceTakingRead:   {},
	AttendanceTakingUpdate: {},
}
