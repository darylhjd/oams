package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
)

func (v *APIServerV1) attendanceTaking(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.attendanceTakingGet(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type attendanceTakingGetResponse struct {
	response
	UpcomingClassGroupSessions []model.ClassGroupSession `json:"upcoming_class_group_sessions"`
}

func (v *APIServerV1) attendanceTakingGet(r *http.Request) apiResponse {
	resp := attendanceTakingGetResponse{
		response:                   newSuccessResponse(),
		UpcomingClassGroupSessions: []model.ClassGroupSession{},
	}

	upcoming, err := v.db.GetUpcomingManagedClassGroupSessions(r.Context())
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get upcoming managed class group sessions")
	}

	resp.UpcomingClassGroupSessions = append(resp.UpcomingClassGroupSessions, upcoming...)
	return resp
}
