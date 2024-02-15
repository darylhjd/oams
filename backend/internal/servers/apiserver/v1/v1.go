package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	dataExportUrl                           = "/data-export"
)

var (
	internalErrorMsg = fmt.Sprintf("%s - internal server error", namespace)
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

	v.mux.HandleFunc(signatureUrl, v.enforceAccessPolicy(
		v.signature,
		map[string][]Permission{
			http.MethodPut: {SignaturePut},
		},
	))

	v.mux.HandleFunc(batchUrl, v.enforceAccessPolicy(
		v.batch,
		map[string][]Permission{
			http.MethodPost: {BatchPost},
			http.MethodPut:  {BatchPut},
		},
	))

	v.mux.HandleFunc(usersUrl, v.enforceAccessPolicy(
		v.users,
		map[string][]Permission{
			http.MethodGet: {UserRead},
		},
	))

	v.mux.HandleFunc(userUrl, v.enforceAccessPolicy(
		v.user,
		map[string][]Permission{
			http.MethodGet:   {UserRead},
			http.MethodPatch: {UserUpdate},
		},
	))

	v.mux.HandleFunc(classesUrl, v.enforceAccessPolicy(
		v.classes,
		map[string][]Permission{
			http.MethodGet: {ClassRead},
		},
	))

	v.mux.HandleFunc(classUrl, v.enforceAccessPolicy(
		v.class,
		map[string][]Permission{
			http.MethodGet: {ClassRead},
		},
	))

	v.mux.HandleFunc(classAttendanceRulesUrl, v.enforceAccessPolicy(
		v.classAttendanceRules,
		map[string][]Permission{
			http.MethodGet: {ClassAttendanceRulesRead},
		},
	))

	v.mux.HandleFunc(classGroupManagersUrl, v.enforceAccessPolicy(
		v.classGroupManagers,
		map[string][]Permission{
			http.MethodGet:  {ClassGroupManagerRead},
			http.MethodPost: {ClassGroupManagerPost},
			http.MethodPut:  {ClassGroupManagerPut},
		},
	))

	v.mux.HandleFunc(classGroupsUrl, v.enforceAccessPolicy(
		v.classGroups,
		map[string][]Permission{
			http.MethodGet: {ClassGroupRead},
		},
	))

	v.mux.HandleFunc(classGroupUrl, v.enforceAccessPolicy(
		v.classGroup,
		map[string][]Permission{
			http.MethodGet: {ClassGroupRead},
		},
	))

	v.mux.HandleFunc(classGroupSessionsUrl, v.enforceAccessPolicy(
		v.classGroupSessions,
		map[string][]Permission{
			http.MethodGet: {ClassGroupSessionRead},
		},
	))

	v.mux.HandleFunc(classGroupSessionUrl, v.enforceAccessPolicy(
		v.classGroupSession,
		map[string][]Permission{
			http.MethodGet: {ClassGroupSessionRead},
		},
	))

	v.mux.HandleFunc(sessionEnrollmentsUrl, v.enforceAccessPolicy(
		v.sessionEnrollments,
		map[string][]Permission{
			http.MethodGet: {SessionEnrollmentRead},
		},
	))

	v.mux.HandleFunc(sessionEnrollmentUrl, v.enforceAccessPolicy(
		v.sessionEnrollment,
		map[string][]Permission{
			http.MethodGet: {SessionEnrollmentRead},
		},
	))

	v.mux.HandleFunc(upcomingClassGroupSessionsUrl, v.enforceAccessPolicy(
		v.upcomingClassGroupSessions,
		map[string][]Permission{
			http.MethodGet: {UpcomingClassGroupSessionRead},
		},
	))

	v.mux.HandleFunc(upcomingClassGroupSessionAttendancesUrl, v.enforceAccessPolicy(
		v.upcomingClassGroupSessionAttendances,
		map[string][]Permission{
			http.MethodGet: {UpcomingClassGroupSessionAttendanceRead},
		},
	))

	v.mux.HandleFunc(upcomingClassGroupSessionAttendanceUrl, v.enforceAccessPolicy(
		v.upcomingClassGroupSessionAttendance,
		map[string][]Permission{
			http.MethodPatch: {UpcomingClassGroupSessionAttendanceUpdate},
		},
	))

	v.mux.HandleFunc(coordinatingClassesUrl, v.enforceAccessPolicy(
		v.coordinatingClasses,
		map[string][]Permission{
			http.MethodGet: {CoordinatingClassRead},
		},
	))

	v.mux.HandleFunc(coordinatingClassUrl, v.enforceAccessPolicy(
		v.coordinatingClass,
		map[string][]Permission{
			http.MethodGet: {CoordinatingClassRead},
		},
	))

	v.mux.HandleFunc(coordinatingClassRulesUrl, v.enforceAccessPolicy(
		v.coordinatingClassRules,
		map[string][]Permission{
			http.MethodGet:  {CoordinatingClassRuleRead},
			http.MethodPost: {CoordinatingClassRuleCreate},
		},
	))

	v.mux.HandleFunc(coordinatingClassRuleUrl, v.enforceAccessPolicy(
		v.coordinatingClassRule,
		map[string][]Permission{
			http.MethodPatch:  {CoordinatingClassRuleUpdate},
			http.MethodDelete: {CoordinatingClassRuleDelete},
		},
	))

	v.mux.HandleFunc(coordinatingClassReportUrl, v.enforceAccessPolicy(
		v.coordinatingClassReport,
		map[string][]Permission{
			http.MethodGet: {CoordinatingClassReportRead},
		},
	))

	v.mux.HandleFunc(coordinatingClassDashboardUrl, v.enforceAccessPolicy(
		v.coordinatingClassDashboard,
		map[string][]Permission{
			http.MethodGet: {CoordinatingClassDashboardRead},
		},
	))

	v.mux.HandleFunc(dataExportUrl, v.enforceAccessPolicy(
		v.dataExport,
		map[string][]Permission{
			http.MethodGet: {DataExportRead},
		},
	))
}

func (v *APIServerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.l.Debug(fmt.Sprintf("%s - handling request", namespace), zap.String("endpoint", r.URL.Path))

	v.mux.ServeHTTP(w, r)
}

func (v *APIServerV1) parseRequestBody(body io.ReadCloser, a any) error {
	var b bytes.Buffer

	if _, err := b.ReadFrom(body); err != nil {
		return err
	}

	return json.Unmarshal(b.Bytes(), a)
}

func (v *APIServerV1) writeResponse(w http.ResponseWriter, r *http.Request, resp apiResponse) {
	b, err := json.Marshal(resp)
	if err != nil {
		v.l.Error(fmt.Sprintf("%s - could not marshal response", namespace),
			zap.String("endpoint", r.URL.Path),
			zap.String("method", r.Method),
			zap.Error(err),
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code())
	if _, err = w.Write(b); err != nil {
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("endpoint", r.URL.Path),
			zap.String("method", r.Method),
			zap.Error(err),
		)
	}
}

func (v *APIServerV1) logInternalServerError(r *http.Request, err error) {
	v.l.Error(internalErrorMsg, zap.String("endpoint", r.URL.Path), zap.String("method", r.Method), zap.Error(err))
}
