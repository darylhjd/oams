package v1

import (
	"net/http"
	"time"

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
	UpcomingClassGroupSessions []attendanceTakingGetUpcomingClassGroupSessionResponseFields `json:"upcoming_class_group_sessions"`
}

type attendanceTakingGetUpcomingClassGroupSessionResponseFields struct {
	ID           int64               `json:"id"`
	StartTime    time.Time           `json:"start_time"`
	EndTime      time.Time           `json:"end_time"`
	Venue        string              `json:"venue"`
	ClassType    model.ClassType     `json:"class_type"`
	ManagingRole *model.ManagingRole `json:"managing_role"` // For nil values, exposed as system admin.
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) attendanceTakingGet(r *http.Request) apiResponse {
	resp := attendanceTakingGetResponse{
		response: newSuccessResponse(),
	}

	upcoming, err := v.db.GetUpcomingManagedClassGroupSessions(r.Context())
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get upcoming managed class group sessions")
	}

	resp.UpcomingClassGroupSessions = make([]attendanceTakingGetUpcomingClassGroupSessionResponseFields, 0, len(upcoming))
	for _, m := range upcoming {
		resp.UpcomingClassGroupSessions = append(resp.UpcomingClassGroupSessions, attendanceTakingGetUpcomingClassGroupSessionResponseFields{
			m.ID,
			m.StartTime,
			m.EndTime,
			m.Venue,
			m.ClassType,
			m.ManagingRole,
		})
	}

	return resp
}
