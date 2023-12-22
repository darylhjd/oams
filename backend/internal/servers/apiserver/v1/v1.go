package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/servers/apiserver/v1/permissions"
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
	baseUrl               = "/"
	pingUrl               = "/ping"
	batchUrl              = "/batch"
	usersUrl              = "/users"
	userUrl               = "/users/"
	classesUrl            = "/classes"
	classUrl              = "/classes/"
	classGroupManagersUrl = "/class-group-managers"
	classGroupsUrl        = "/class-groups"
	classGroupUrl         = "/class-groups/"
	classGroupSessionsUrl = "/class-group-sessions"
	classGroupSessionUrl  = "/class-group-sessions/"
	sessionEnrollmentsUrl = "/session-enrollments"
	sessionEnrollmentUrl  = "/session-enrollments/"
	attendanceTakingsUrl  = "/attendance-taking"
	attendanceTakingUrl   = "/attendance-taking/"
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

	v.mux.HandleFunc(batchUrl, permissions.EnforceAccessPolicy(
		v.batch,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodPost: {permissions.BatchPost},
			http.MethodPut:  {permissions.BatchPut},
		},
	))

	v.mux.HandleFunc(usersUrl, permissions.EnforceAccessPolicy(
		v.users,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.UserRead},
		},
	))

	v.mux.HandleFunc(userUrl, permissions.EnforceAccessPolicy(
		v.user,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet:   {permissions.UserRead},
			http.MethodPatch: {permissions.UserUpdate},
		},
	))

	v.mux.HandleFunc(classesUrl, permissions.EnforceAccessPolicy(
		v.classes,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.ClassRead},
		},
	))

	v.mux.HandleFunc(classUrl, permissions.EnforceAccessPolicy(
		v.class,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.ClassRead},
		},
	))

	v.mux.HandleFunc(classGroupManagersUrl, permissions.EnforceAccessPolicy(
		v.classGroupManagers,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet:  {permissions.ClassManagerRead},
			http.MethodPost: {permissions.ClassManagerCreate},
			http.MethodPut:  {permissions.ClassManagerCreate},
		},
	))

	v.mux.HandleFunc(classGroupsUrl, permissions.EnforceAccessPolicy(
		v.classGroups,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.ClassGroupRead},
		},
	))

	v.mux.HandleFunc(classGroupUrl, permissions.EnforceAccessPolicy(
		v.classGroup,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.ClassGroupRead},
		},
	))

	v.mux.HandleFunc(classGroupSessionsUrl, permissions.EnforceAccessPolicy(
		v.classGroupSessions,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.ClassGroupSessionRead},
		},
	))

	v.mux.HandleFunc(classGroupSessionUrl, permissions.EnforceAccessPolicy(
		v.classGroupSession,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.ClassGroupSessionRead},
		},
	))

	v.mux.HandleFunc(sessionEnrollmentsUrl, permissions.EnforceAccessPolicy(
		v.sessionEnrollments,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.SessionEnrollmentRead},
		},
	))

	v.mux.HandleFunc(sessionEnrollmentUrl, permissions.EnforceAccessPolicy(
		v.sessionEnrollment,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet:   {permissions.SessionEnrollmentRead},
			http.MethodPatch: {permissions.SessionEnrollmentUpdate},
		},
	))

	v.mux.HandleFunc(attendanceTakingsUrl, permissions.EnforceAccessPolicy(
		v.attendanceTakings,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.AttendanceTakingRead},
		},
	))

	v.mux.HandleFunc(attendanceTakingUrl, permissions.EnforceAccessPolicy(
		v.attendanceTaking,
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.AttendanceTakingRead},
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
