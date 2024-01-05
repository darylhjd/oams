package v1

import (
	"fmt"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/v1/permissions"
)

const (
	coordinatingClassRulesUrl  = "/rules"
	coordinatingClassReportUrl = "/report"
)

var (
	coordinatingClassRulesFormat = fmt.Sprintf("%s%%d/%%s", coordinatingClassUrl)
)

func (v *APIServerV1) coordinatingClass(w http.ResponseWriter, r *http.Request) {
	var (
		classId int64
		throw   string
	)
	if _, err := fmt.Sscanf(r.URL.Path, coordinatingClassRulesFormat, &classId, &throw); err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, ""))
		return
	}

	mux := http.NewServeMux()

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

	prefix := fmt.Sprintf("%s%d", coordinatingClassUrl, classId)
	http.StripPrefix(prefix, mux).ServeHTTP(w, r)
}
