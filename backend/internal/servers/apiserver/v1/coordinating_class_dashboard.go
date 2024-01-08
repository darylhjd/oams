package v1

import (
	"net/http"
)

func (v *APIServerV1) coordinatingClassDashboard(w http.ResponseWriter, r *http.Request, classId int64) {
	var resp apiResponse

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
}

func (v *APIServerV1) coordinatingClassDashboardGet(_ *http.Request, _ int64) apiResponse {
	return coordinatingClassDashboardGetResponse{
		newSuccessResponse(),
	}
}
