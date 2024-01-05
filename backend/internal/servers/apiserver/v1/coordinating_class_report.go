package v1

import "net/http"

func (v *APIServerV1) coordinatingClassReport(w http.ResponseWriter, r *http.Request, classId int64) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.coordinatingClassReportGet(r, classId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

func (v *APIServerV1) coordinatingClassReportGet(r *http.Request, id int64) apiResponse {
	return newErrorResponse(http.StatusNotImplemented, "")
}
