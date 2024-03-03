package v1

import (
	"net/http"
	"slices"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

type permission int

const (
	SignaturePut permission = iota

	BatchPost
	BatchPut

	UserRead
	UserUpdate

	ClassRead

	ClassAttendanceRulesRead

	ClassGroupManagerPost
	ClassGroupManagerRead
	ClassGroupManagerUpdate
	ClassGroupManagerDelete
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

type permissionMap map[permission]struct{}

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

	ClassGroupManagerPost:   {},
	ClassGroupManagerRead:   {},
	ClassGroupManagerUpdate: {},
	ClassGroupManagerDelete: {},
	ClassGroupManagerPut:    {},

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
func hasPermissions(role model.UserRole, p permission) bool {
	permModel, ok := rolePermissionMapping[role]
	if !ok {
		return false
	}

	_, ok = permModel[p]
	return ok
}

func (v *APIServerV1) enforcePermissionAccess(
	handlerFunc http.HandlerFunc,
	methodPermissions map[string]permission,
) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		authContext := oauth2.GetAuthContext(r.Context())

		if !hasPermissions(authContext.User.Role, methodPermissions[r.Method]) {
			v.writeResponse(w, r, newErrorResponse(http.StatusUnauthorized, "insufficient permissions"))
			return
		}

		handlerFunc(w, r)
	}

	return middleware.MustAuth(handler, v.auth, v.db)
}

const (
	roleAttendanceTaker = "Attendance.Take"
)

func (v *APIServerV1) enforceRoleAccess(
	handlerFunc http.HandlerFunc,
	allowedRoles []string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roles := oauth2.GetAuthContext(r.Context()).Claims.AppRoles()

		for _, role := range roles {
			if !slices.Contains(allowedRoles, role) {
				v.writeResponse(w, r, newErrorResponse(http.StatusUnauthorized, "insufficient scope"))
				return
			}
		}

		handlerFunc(w, r)
	}
}

// enforceAccess is a helper function for access control. Depending on the type of user, the relevant access policy is
// enforced - permissions for users and role for applications.
func (v *APIServerV1) enforceAccess(
	handlerFunc http.HandlerFunc,
	methodPermissions map[string]permission,
	allowedRoles []string,
) http.HandlerFunc {
	return middleware.MustAuth(func(w http.ResponseWriter, r *http.Request) {
		authContext := oauth2.GetAuthContext(r.Context())

		switch {
		case authContext.Claims.IsApplication():
			v.enforceRoleAccess(handlerFunc, allowedRoles)(w, r)
		default:
			v.enforcePermissionAccess(handlerFunc, methodPermissions)(w, r)
		}
	}, v.auth, v.db)
}
