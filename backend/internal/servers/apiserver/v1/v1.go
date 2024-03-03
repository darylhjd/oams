package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

const namespace = "apiserver/v1"

const (
	Url = "/api/v1/"
)

const (
	baseUrl                                 = "/"
	pingUrl                                 = "/ping"
	loginUrl                                = "/login"
	msLoginCallbackUrl                      = "/ms-login-callback"
	logoutUrl                               = "/logout"
	sessionUrl                              = "/session"
	signatureUrl                            = "/signature/{userId}"
	batchUrl                                = "/batch"
	usersUrl                                = "/users"
	userUrl                                 = "/users/{userId}"
	classesUrl                              = "/classes"
	classUrl                                = "/classes/{classId}"
	classAttendanceRulesUrl                 = "/class-attendance-rules"
	classGroupManagersUrl                   = "/class-group-managers"
	classGroupManagerUrl                    = "/class-group-managers/{managerId}"
	classGroupsUrl                          = "/class-groups"
	classGroupUrl                           = "/class-groups/{groupId}"
	classGroupSessionsUrl                   = "/class-group-sessions"
	classGroupSessionUrl                    = "/class-group-sessions/{sessionId}"
	sessionEnrollmentsUrl                   = "/session-enrollments"
	sessionEnrollmentUrl                    = "/session-enrollments/{enrollmentId}"
	upcomingClassGroupSessionsUrl           = "/upcoming-class-group-sessions"
	upcomingClassGroupSessionAttendancesUrl = "/upcoming-class-group-sessions/{sessionId}/attendances"
	upcomingClassGroupSessionAttendanceUrl  = "/upcoming-class-group-sessions/{sessionId}/attendances/{enrollmentId}"
	coordinatingClassesUrl                  = "/coordinating-classes"
	coordinatingClassUrl                    = "/coordinating-classes/{classId}"
	coordinatingClassRulesUrl               = "/coordinating-classes/{classId}/rules"
	coordinatingClassRuleUrl                = "/coordinating-classes/{classId}/rules/{ruleId}"
	coordinatingClassReportUrl              = "/coordinating-classes/{classId}/report"
	coordinatingClassDashboardUrl           = "/coordinating-classes/{classId}/dashboard"
	coordinatingClassSchedulesUrl           = "/coordinating-classes/{classId}/schedule"
	coordinatingClassScheduleUrl            = "/coordinating-classes/{classId}/schedule/{sessionId}"
	dataExportUrl                           = "/data-export"
)

type APIServerV1 struct {
	l   *zap.Logger
	db  *database.DB
	mux *http.ServeMux

	decoder *schema.Decoder

	auth oauth2.AuthProvider
}

// New creates a new APIServerV1. This is a sub-router and should not be used as a base router.
func New(l *zap.Logger, db *database.DB, auth oauth2.AuthProvider) *APIServerV1 {
	server := APIServerV1{l, db, http.NewServeMux(), schema.NewDecoder(), auth}
	server.registerHandlers()

	return &server
}

func (v *APIServerV1) registerHandlers() {
	v.mux.HandleFunc(baseUrl, v.base)
	v.mux.HandleFunc(pingUrl, v.ping)

	v.mux.HandleFunc(loginUrl, v.login)
	v.mux.HandleFunc(msLoginCallbackUrl, v.msLoginCallback)
	v.mux.HandleFunc(logoutUrl, v.logout)

	v.mux.HandleFunc(sessionUrl, v.session)

	v.mux.HandleFunc(signatureUrl, v.enforceAccess(
		v.signature,
		map[string]permission{
			http.MethodPut: SignaturePut,
		},
		[]string{},
	))

	v.mux.HandleFunc(batchUrl, v.enforceAccess(
		v.batch,
		map[string]permission{
			http.MethodPost: BatchPost,
			http.MethodPut:  BatchPut,
		},
		[]string{},
	))

	v.mux.HandleFunc(usersUrl, v.enforceAccess(
		v.users,
		map[string]permission{
			http.MethodGet: UserRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(userUrl, v.enforceAccess(
		v.user,
		map[string]permission{
			http.MethodGet:   UserRead,
			http.MethodPatch: UserUpdate,
		},
		[]string{},
	))

	v.mux.HandleFunc(classesUrl, v.enforceAccess(
		v.classes,
		map[string]permission{
			http.MethodGet: ClassRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(classUrl, v.enforceAccess(
		v.class,
		map[string]permission{
			http.MethodGet: ClassRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(classAttendanceRulesUrl, v.enforceAccess(
		v.classAttendanceRules,
		map[string]permission{
			http.MethodGet: ClassAttendanceRulesRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(classGroupManagersUrl, v.enforceAccess(
		v.classGroupManagers,
		map[string]permission{
			http.MethodGet:  ClassGroupManagerRead,
			http.MethodPost: ClassGroupManagerPost,
			http.MethodPut:  ClassGroupManagerPut,
		},
		[]string{},
	))

	v.mux.HandleFunc(classGroupManagerUrl, v.enforceAccess(
		v.classGroupManager,
		map[string]permission{
			http.MethodGet:    ClassGroupManagerRead,
			http.MethodPatch:  ClassGroupManagerUpdate,
			http.MethodDelete: ClassGroupManagerDelete,
		},
		[]string{},
	))

	v.mux.HandleFunc(classGroupsUrl, v.enforceAccess(
		v.classGroups,
		map[string]permission{
			http.MethodGet: ClassGroupRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(classGroupUrl, v.enforceAccess(
		v.classGroup,
		map[string]permission{
			http.MethodGet: ClassGroupRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(classGroupSessionsUrl, v.enforceAccess(
		v.classGroupSessions,
		map[string]permission{
			http.MethodGet: ClassGroupSessionRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(classGroupSessionUrl, v.enforceAccess(
		v.classGroupSession,
		map[string]permission{
			http.MethodGet: ClassGroupSessionRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(sessionEnrollmentsUrl, v.enforceAccess(
		v.sessionEnrollments,
		map[string]permission{
			http.MethodGet: SessionEnrollmentRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(sessionEnrollmentUrl, v.enforceAccess(
		v.sessionEnrollment,
		map[string]permission{
			http.MethodGet: SessionEnrollmentRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(upcomingClassGroupSessionsUrl, v.enforceAccess(
		v.upcomingClassGroupSessions,
		map[string]permission{
			http.MethodGet: UpcomingClassGroupSessionRead,
		},
		[]string{roleAttendanceTaker},
	))

	v.mux.HandleFunc(upcomingClassGroupSessionAttendancesUrl, v.enforceAccess(
		v.upcomingClassGroupSessionAttendances,
		map[string]permission{
			http.MethodGet: UpcomingClassGroupSessionAttendanceRead,
		},
		[]string{roleAttendanceTaker},
	))

	v.mux.HandleFunc(upcomingClassGroupSessionAttendanceUrl, v.enforceAccess(
		v.upcomingClassGroupSessionAttendance,
		map[string]permission{
			http.MethodPatch: UpcomingClassGroupSessionAttendanceUpdate,
		},
		[]string{roleAttendanceTaker},
	))

	v.mux.HandleFunc(coordinatingClassesUrl, v.enforceAccess(
		v.coordinatingClasses,
		map[string]permission{
			http.MethodGet: CoordinatingClassRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(coordinatingClassUrl, v.enforceAccess(
		v.coordinatingClass,
		map[string]permission{
			http.MethodGet: CoordinatingClassRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(coordinatingClassRulesUrl, v.enforceAccess(
		v.coordinatingClassRules,
		map[string]permission{
			http.MethodGet:  CoordinatingClassRuleRead,
			http.MethodPost: CoordinatingClassRuleCreate,
		},
		[]string{},
	))

	v.mux.HandleFunc(coordinatingClassRuleUrl, v.enforceAccess(
		v.coordinatingClassRule,
		map[string]permission{
			http.MethodPatch:  CoordinatingClassRuleUpdate,
			http.MethodDelete: CoordinatingClassRuleDelete,
		},
		[]string{},
	))

	v.mux.HandleFunc(coordinatingClassReportUrl, v.enforceAccess(
		v.coordinatingClassReport,
		map[string]permission{
			http.MethodGet: CoordinatingClassReportRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(coordinatingClassDashboardUrl, v.enforceAccess(
		v.coordinatingClassDashboard,
		map[string]permission{
			http.MethodGet: CoordinatingClassDashboardRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(coordinatingClassSchedulesUrl, v.enforceAccess(
		v.coordinatingClassSchedules,
		map[string]permission{
			http.MethodGet: CoordinatingClassScheduleRead,
		},
		[]string{},
	))

	v.mux.HandleFunc(coordinatingClassScheduleUrl, v.enforceAccess(
		v.coordinatingClassSchedule,
		map[string]permission{
			http.MethodPut: CoordinatingClassScheduleUpdate,
		},
		[]string{},
	))

	v.mux.HandleFunc(dataExportUrl, v.enforceAccess(
		v.dataExport,
		map[string]permission{
			http.MethodGet: DataExportRead,
		},
		[]string{},
	))
}

func (v *APIServerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.l.Debug(fmt.Sprintf("%s - handling request", namespace), zap.String("endpoint", r.URL.Path))
	v.mux.ServeHTTP(w, r)
}

func (v *APIServerV1) writeResponse(w http.ResponseWriter, r *http.Request, resp apiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code())
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("endpoint", r.URL.Path),
			zap.String("method", r.Method),
			zap.Error(err),
		)
	}
}

func (v *APIServerV1) logInternalServerError(r *http.Request, err error) {
	v.l.Error(
		fmt.Sprintf("%s - internal server error", namespace),
		zap.String("endpoint", r.URL.Path), zap.String("method", r.Method), zap.Error(err),
	)
}
