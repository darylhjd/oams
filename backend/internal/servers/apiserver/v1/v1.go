package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/permissions"
	"github.com/gorilla/schema"
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

const namespace = "apiserver/v1"

const (
	Url = "/api/v1/"
)

const (
	baseUrl               = "/"
	pingUrl               = "/ping"
	loginUrl              = "/login"
	msLoginCallbackUrl    = "/ms-login-callback"
	logoutUrl             = "/logout"
	batchUrl              = "/batch"
	usersUrl              = "/users"
	userUrl               = "/users/"
	classesUrl            = "/classes"
	classUrl              = "/classes/"
	classGroupsUrl        = "/class-groups"
	classGroupUrl         = "/class-groups/"
	classGroupSessionsUrl = "/class-group-sessions"
	classGroupSessionUrl  = "/class-group-sessions/"
	sessionEnrollmentsUrl = "/session-enrollments"
	sessionEnrollmentUrl  = "/session-enrollments/"
)

var (
	internalErrorMsg = fmt.Sprintf("%s - internal server error", namespace)
)

type APIServerV1 struct {
	l   *zap.Logger
	db  *database.DB
	mux *http.ServeMux

	decoder *schema.Decoder

	azure oauth2.Authenticator
}

// New creates a new APIServerV1. This is a sub-router and should not be used as a base router.
func New(l *zap.Logger, db *database.DB, azureClient oauth2.Authenticator) *APIServerV1 {
	server := APIServerV1{l, db, http.NewServeMux(), schema.NewDecoder(), azureClient}
	server.registerHandlers()

	return &server
}

func (v *APIServerV1) registerHandlers() {
	v.mux.HandleFunc(baseUrl, middleware.AllowMethods(v.base, http.MethodGet))
	v.mux.HandleFunc(pingUrl, middleware.AllowMethods(v.ping, http.MethodGet))

	v.mux.HandleFunc(loginUrl, v.login)
	v.mux.HandleFunc(msLoginCallbackUrl, middleware.AllowMethods(v.msLoginCallback, http.MethodPost))
	v.mux.HandleFunc(logoutUrl, middleware.MustAuth(v.logout, v.azure, v.db))

	v.mux.HandleFunc(batchUrl, middleware.MustAuth(
		middleware.AllowMethodsWithPermissions(v.batch,
			map[string][]permissions.Permission{
				http.MethodPost: {permissions.PermissionBatchPost},
				http.MethodPut:  {permissions.PermissionBatchPut},
			},
		),
		v.azure, v.db,
	))

	v.mux.HandleFunc(usersUrl, middleware.MustAuth(
		middleware.AllowMethodsWithPermissions(v.users,
			map[string][]permissions.Permission{
				http.MethodGet:  {permissions.PermissionUserRead},
				http.MethodPost: {permissions.PermissionUserCreate},
			},
		),
		v.azure, v.db,
	))

	v.mux.HandleFunc(userUrl, middleware.MustAuth(
		middleware.AllowMethods(v.user, http.MethodGet, http.MethodPatch, http.MethodDelete),
		v.azure, v.db,
	))

	v.mux.HandleFunc(classesUrl, middleware.MustAuth(
		middleware.AllowMethods(v.classes, http.MethodGet, http.MethodPost),
		v.azure, v.db,
	))

	v.mux.HandleFunc(classUrl, middleware.MustAuth(
		middleware.AllowMethods(v.class, http.MethodGet, http.MethodPatch, http.MethodDelete),
		v.azure, v.db,
	))

	v.mux.HandleFunc(classGroupsUrl, middleware.MustAuth(
		middleware.AllowMethods(v.classGroups, http.MethodGet, http.MethodPost),
		v.azure, v.db,
	))

	v.mux.HandleFunc(classGroupUrl, middleware.MustAuth(
		middleware.AllowMethods(v.classGroup, http.MethodGet, http.MethodPatch, http.MethodDelete),
		v.azure, v.db,
	))

	v.mux.HandleFunc(classGroupSessionsUrl, middleware.MustAuth(
		middleware.AllowMethods(v.classGroupSessions, http.MethodGet, http.MethodPost),
		v.azure, v.db,
	))

	v.mux.HandleFunc(classGroupSessionUrl, middleware.MustAuth(
		middleware.AllowMethods(v.classGroupSession, http.MethodGet, http.MethodPatch, http.MethodDelete),
		v.azure, v.db,
	))

	v.mux.HandleFunc(sessionEnrollmentsUrl, middleware.MustAuth(
		middleware.AllowMethods(v.sessionEnrollments, http.MethodGet, http.MethodPost),
		v.azure, v.db,
	))

	v.mux.HandleFunc(sessionEnrollmentUrl, middleware.MustAuth(
		middleware.AllowMethods(v.sessionEnrollment, http.MethodGet, http.MethodPatch, http.MethodDelete),
		v.azure, v.db,
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
