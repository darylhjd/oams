package v1

import (
	"net/http"
	"strconv"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) coordinatingClassDashboard(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classId, err := strconv.ParseInt(r.PathValue("classId"), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.coordinatingClassDashboardGet(r, classId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type coordinatingClassDashboardGetResponse struct {
	response
	AttendanceCount []database.AttendanceCountData `json:"attendance_count"`
}

func (v *APIServerV1) coordinatingClassDashboardGet(r *http.Request, classId int64) apiResponse {
	data, err := v.db.GetDashboardData(r.Context(), classId)
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, "could not get dashboard data")
	}

	return coordinatingClassDashboardGetResponse{
		newSuccessResponse(),
		append(make([]database.AttendanceCountData, 0, len(data)), data...),
	}
}
