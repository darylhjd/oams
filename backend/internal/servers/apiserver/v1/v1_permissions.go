package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

type Permission int

const (
	SignaturePut Permission = iota

	BatchPost
	BatchPut

	UserRead
	UserUpdate

	ClassRead

	ClassAttendanceRulesRead

	ClassGroupManagerPost
	ClassGroupManagerRead
	ClassGroupManagerPut

	ClassGroupRead

	ClassGroupSessionRead

	SessionEnrollmentRead

	UpcomingClassGroupSessionRead

	UpcomingClassGroupSessionAttendanceRead
	UpcomingClassGroupSessionAttendanceUpdate

	CoordinatingClassRead

	CoordinatingClassRuleCreate
	CoordinatingClassRuleRead
	CoordinatingClassRuleUpdate
	CoordinatingClassRuleDelete

	CoordinatingClassReportRead

	CoordinatingClassDashboardRead

	CoordinatingClassScheduleRead
	CoordinatingClassScheduleUpdate

	DataExportRead
)

type permissionMap map[Permission]struct{}

var rolePermissionMapping = map[model.UserRole]permissionMap{
	model.UserRole_User:        userRolePermissions,
	model.UserRole_SystemAdmin: systemAdminRolePermissions,
}

var userRolePermissions = permissionMap{
	SignaturePut: {},

	UserRead:   {},
	UserUpdate: {},

	ClassRead: {},

	ClassGroupRead: {},

	ClassGroupSessionRead: {},

	SessionEnrollmentRead: {},

	UpcomingClassGroupSessionRead: {},

	UpcomingClassGroupSessionAttendanceRead:   {},
	UpcomingClassGroupSessionAttendanceUpdate: {},

	CoordinatingClassRead: {},

	CoordinatingClassRuleCreate: {},
	CoordinatingClassRuleRead:   {},
	CoordinatingClassRuleUpdate: {},
	CoordinatingClassRuleDelete: {},

	CoordinatingClassDashboardRead: {},

	CoordinatingClassReportRead: {},

	CoordinatingClassScheduleRead:   {},
	CoordinatingClassScheduleUpdate: {},
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

	UpcomingClassGroupSessionRead: {},

	UpcomingClassGroupSessionAttendanceRead:   {},
	UpcomingClassGroupSessionAttendanceUpdate: {},

	CoordinatingClassRead: {},

	CoordinatingClassRuleCreate: {},
	CoordinatingClassRuleRead:   {},
	CoordinatingClassRuleUpdate: {},
	CoordinatingClassRuleDelete: {},

	CoordinatingClassReportRead: {},

	CoordinatingClassDashboardRead: {},

	CoordinatingClassScheduleRead:   {},
	CoordinatingClassScheduleUpdate: {},

	DataExportRead: {},
}

// hasPermissions checks if a user with a role has all the given permissions.
func hasPermissions(role model.UserRole, permissions ...Permission) bool {
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

// enforceAccessPolicy based on role-based access control.
func (v *APIServerV1) enforceAccessPolicy(
	handlerFunc http.HandlerFunc,
	methodPermissions map[string][]Permission,
) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		authContext := oauth2.GetAuthContext(r.Context())

		if !hasPermissions(authContext.User.Role, methodPermissions[r.Method]...) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handlerFunc(w, r)
	}

	return middleware.MustAuth(handler, v.auth, v.db)
}
