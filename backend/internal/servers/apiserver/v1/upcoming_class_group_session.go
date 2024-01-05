package v1

import (
	"fmt"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/v1/permissions"
)

const (
	upcomingClassGroupSessionAttendancesUrl = "/attendances"
	upcomingClassGroupSessionAttendanceUrl  = "/attendances/"
)

var (
	upcomingClassGroupSessionSubFormat = fmt.Sprintf("%s%%d/%%s", upcomingClassGroupSessionUrl)
)

func (v *APIServerV1) upcomingClassGroupSession(w http.ResponseWriter, r *http.Request) {
	var (
		sessionId int64
		throw     string
	)
	if _, err := fmt.Sscanf(r.URL.Path, upcomingClassGroupSessionSubFormat, &sessionId, &throw); err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group session id"))
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc(upcomingClassGroupSessionAttendancesUrl, permissions.EnforceAccessPolicy(
		middleware.WithID(sessionId, v.upcomingClassGroupSessionAttendances),
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodGet: {permissions.UpcomingClassGroupSessionAttendanceRead},
		},
	))

	mux.HandleFunc(upcomingClassGroupSessionAttendanceUrl, permissions.EnforceAccessPolicy(
		middleware.WithID(sessionId, v.upcomingClassGroupSessionAttendance),
		v.auth, v.db,
		map[string][]permissions.P{
			http.MethodPatch: {permissions.UpcomingClassGroupSessionAttendanceUpdate},
		},
	))

	prefix := fmt.Sprintf("%s%d", upcomingClassGroupSessionUrl, sessionId)
	http.StripPrefix(prefix, mux).ServeHTTP(w, r)
}
