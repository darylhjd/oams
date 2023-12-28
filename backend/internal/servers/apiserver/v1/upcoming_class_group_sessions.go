package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) upcomingClassGroupSessions(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.upcomingClassGroupSessionsGet(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type upcomingClassGroupSessionsGetResponse struct {
	response
	UpcomingClassGroupSessions []database.UpcomingManagedClassGroupSession `json:"upcoming_class_group_sessions"`
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) upcomingClassGroupSessionsGet(r *http.Request) apiResponse {
	upcoming, err := v.db.GetUpcomingManagedClassGroupSessions(r.Context())
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get upcoming class group sessions")
	}

	return upcomingClassGroupSessionsGetResponse{
		newSuccessResponse(),
		append(
			make([]database.UpcomingManagedClassGroupSession, 0, len(upcoming)),
			upcoming...,
		),
	}
}
