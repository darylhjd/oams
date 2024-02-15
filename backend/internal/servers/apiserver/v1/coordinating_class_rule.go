package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) coordinatingClassRule(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classId, err := strconv.ParseInt(r.PathValue("classId"), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	ruleId, err := strconv.ParseInt(r.PathValue("ruleId"), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid rule id"))
		return
	}

	switch r.Method {
	case http.MethodPatch:
		resp = v.coordinatingClassRulePatch(r, classId, ruleId)
	case http.MethodDelete:
		resp = v.coordinatingClassRuleDelete(r, classId, ruleId)
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

func (v *APIServerV1) coordinatingClassRulePatch(r *http.Request, classId, ruleId int64) apiResponse {
	var req coordinatingClassRulePatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	active, err := v.db.UpdateCoordinatingClassRule(r.Context(), classId, ruleId, req.Active)
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

type coordinatingClassRuleDeleteResponse struct {
	response
}

func (v *APIServerV1) coordinatingClassRuleDelete(r *http.Request, classId, ruleId int64) apiResponse {
	if err := v.db.DeleteCoordinatingClassRule(r.Context(), classId, ruleId); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusUnauthorized, "not allowed to delete rule")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process coordinating class rule delete database action")
	}

	return coordinatingClassRuleDeleteResponse{
		newSuccessResponse(),
	}
}
