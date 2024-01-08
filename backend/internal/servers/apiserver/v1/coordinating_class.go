package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/v1/permissions"
	"github.com/go-jet/jet/v2/qrm"
)

const (
	coordinatingClassBaseUrl      = "/"
	coordinatingClassRulesUrl     = "/rules"
	coordinatingClassReportUrl    = "/report"
	coordinatingClassDashboardUrl = "/dashboard"
)

type coordinatingClassMux struct {
	mux *http.ServeMux
}

func (c *coordinatingClassMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}

func (v *APIServerV1) newCoordinatingClassMux(classId int64) *coordinatingClassMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", permissions.EnforceAccessPolicy(
		middleware.WithID(classId, v.coordinatingClassBase),
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.CoordinatingClassRead},
		},
	))

	mux.HandleFunc(coordinatingClassRulesUrl, permissions.EnforceAccessPolicy(
		middleware.WithID(classId, v.coordinatingClassRules),
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet:  {permissions.CoordinatingClassRuleRead},
			http.MethodPost: {permissions.CoordinatingClassRuleCreate},
		},
	))

	mux.HandleFunc(coordinatingClassReportUrl, permissions.EnforceAccessPolicy(
		middleware.WithID(classId, v.coordinatingClassReport),
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.CoordinatingClassReportRead},
		},
	))

	mux.HandleFunc(coordinatingClassDashboardUrl, permissions.EnforceAccessPolicy(
		middleware.WithID(classId, v.coordinatingClassDashboard),
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.CoordinatingClassDashboardRead},
		},
	))

	return &coordinatingClassMux{mux}
}

func (v *APIServerV1) coordinatingClass(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, coordinatingClassUrl), "/", 2)
	if len(parts) < 1 {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	classId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	prefix := fmt.Sprintf("%s%d", coordinatingClassUrl, classId)
	if strings.TrimPrefix(prefix, r.URL.Path) == "" { // Help with base URL.
		r.URL.Path += coordinatingClassBaseUrl
	}

	http.StripPrefix(prefix, v.newCoordinatingClassMux(classId)).ServeHTTP(w, r)
}

func (v *APIServerV1) coordinatingClassBase(w http.ResponseWriter, r *http.Request, classId int64) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.coordinatingClassGet(r, classId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type coordinatingClassGetResponse struct {
	response
	CoordinatingClass database.CoordinatingClass `json:"coordinating_class"`
}

func (v *APIServerV1) coordinatingClassGet(r *http.Request, classId int64) apiResponse {
	class, err := v.db.GetCoordinatingClass(r.Context(), classId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested coordinating class does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process coordinating class get database action")
	}

	return coordinatingClassGetResponse{
		newSuccessResponse(),
		class,
	}
}
