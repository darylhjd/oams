package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/pkg/to"
)

func (v *APIServerV1) coordinatingClassSchedules(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classId, err := to.Int64(r.PathValue("classId"))
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.coordinatingClassSchedulesGet(r, classId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type coordinatingClassSchedulesGetResponse struct {
	response
	Schedule []database.ScheduleData `json:"schedule"`
}

func (v *APIServerV1) coordinatingClassSchedulesGet(r *http.Request, classId int64) apiResponse {
	schedule, err := v.db.GetCoordinatingClassSchedules(r.Context(), classId)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get coordinating class schedule")
	}

	return coordinatingClassSchedulesGetResponse{
		newSuccessResponse(),
		schedule,
	}
}
