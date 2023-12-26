package v1

import "net/http"

func (v *APIServerV1) attendanceRules(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.attendanceRulesGet(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type attendanceRulesGetResponse struct {
	response
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) attendanceRulesGet(_ *http.Request) apiResponse {
	return attendanceRulesGetResponse{newSuccessResponse()}
}
