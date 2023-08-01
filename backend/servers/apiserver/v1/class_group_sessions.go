package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) classGroupSessions(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupSessionsGet(r)
	case http.MethodPost:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, classGroupSessionsUrl, resp)
}

type classGroupSessionsGetResponse struct {
	response
	ClassGroupSessions []database.ClassGroupSession `json:"class_group_sessions"`
}

func (v *APIServerV1) classGroupSessionsGet(r *http.Request) apiResponse {
	sessions, err := v.db.Q.ListClassGroupSessions(r.Context())
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, "could not process class group sessions get database action")
	}

	resp := classGroupSessionsGetResponse{
		newSuccessResponse(),
		make([]database.ClassGroupSession, 0, len(sessions)),
	}

	resp.ClassGroupSessions = append(resp.ClassGroupSessions, sessions...)
	return resp
}
