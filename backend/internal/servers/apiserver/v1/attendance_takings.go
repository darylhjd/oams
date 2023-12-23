package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) attendanceTakings(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.attendanceTakingsGet(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type attendanceTakingsGetResponse struct {
	response
	UpcomingClassGroupSessions []database.UpcomingManagedClassGroupSession `json:"upcoming_class_group_sessions"`
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) attendanceTakingsGet(r *http.Request) apiResponse {
	resp := attendanceTakingsGetResponse{
		response: newSuccessResponse(),
	}

	upcoming, err := v.db.GetUpcomingManagedClassGroupSessions(r.Context())
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get upcoming managed class group sessions")
	}

	resp.UpcomingClassGroupSessions = append(
		make([]database.UpcomingManagedClassGroupSession, 0, len(upcoming)),
		upcoming...,
	)
	return resp
}
