package v1

import (
	"net/http"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
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
	UpcomingClassGroupSessions []upcomingClassGroupSession `json:"upcoming_class_group_sessions"`
}

type upcomingClassGroupSession struct {
	ID           int64               `json:"id"`
	StartTime    time.Time           `json:"start_time"`
	EndTime      time.Time           `json:"end_time"`
	Venue        string              `json:"venue"`
	Code         string              `json:"code"`
	Year         int32               `json:"year"`
	Semester     string              `json:"semester"`
	Name         string              `json:"name"`
	ClassType    model.ClassType     `json:"class_type"`
	ManagingRole *model.ManagingRole `json:"managing_role"`
}

func (u upcomingClassGroupSession) fromDatabaseUpcomingClassGroupSession(s database.UpcomingManagedClassGroupSession) upcomingClassGroupSession {
	return upcomingClassGroupSession{
		s.ID,
		s.StartTime,
		s.EndTime,
		s.Venue,
		s.Code,
		s.Year,
		s.Semester,
		s.Name,
		s.ClassType,
		s.ManagingRole,
	}
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

	resp.UpcomingClassGroupSessions = make([]upcomingClassGroupSession, 0, len(upcoming))
	for _, m := range upcoming {
		resp.UpcomingClassGroupSessions = append(
			resp.UpcomingClassGroupSessions,
			upcomingClassGroupSession{}.fromDatabaseUpcomingClassGroupSession(m),
		)
	}

	return resp
}
