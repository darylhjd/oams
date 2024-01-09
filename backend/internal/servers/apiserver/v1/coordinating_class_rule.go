package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) coordinatingClassRule(w http.ResponseWriter, r *http.Request, _ int64) {
	var resp apiResponse

	ruleId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, coordinatingClassRuleUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid rule id"))
		return
	}

	switch r.Method {
	case http.MethodPatch:
		resp = v.coordinatingClassRulePatch(r, ruleId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type coordinatingClassRulePatchRequest struct {
	Active bool `json:"active"`
}

type coordinatingClassRulePatchResponse struct {
	response
	Active bool `json:"active"`
}

func (v *APIServerV1) coordinatingClassRulePatch(r *http.Request, ruleId int64) apiResponse {
	var req coordinatingClassRulePatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	active, err := v.db.UpdateCoordinatingClassRule(r.Context(), ruleId, req.Active)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusUnauthorized, "not allowed to toggle rule active")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process coordinating class rule patch database action")
	}

	return coordinatingClassRulePatchResponse{
		newSuccessResponse(),
		active,
	}
}
