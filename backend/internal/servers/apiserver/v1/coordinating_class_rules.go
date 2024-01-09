package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	"github.com/darylhjd/oams/backend/internal/rules"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) coordinatingClassRules(w http.ResponseWriter, r *http.Request, classId int64) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.coordinatingClassRulesGet(r, classId)
	case http.MethodPost:
		resp = v.coordinatingClassRulesPost(r, classId)
	case http.MethodPatch:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type coordinatingClassRulesGetResponse struct {
	response
	Rules []model.ClassAttendanceRule `json:"rules"`
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) coordinatingClassRulesGet(r *http.Request, classId int64) apiResponse {
	classRules, err := v.db.GetCoordinatingClassRules(r.Context(), classId)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get coordinating class rules")
	}

	return coordinatingClassRulesGetResponse{
		newSuccessResponse(),
		append(make([]model.ClassAttendanceRule, 0, len(classRules)), classRules...),
	}
}

type coordinatingClassRulesPostRequest struct {
	rules.RuleParams
}

type coordinatingClassRulesPostResponse struct {
	response
	Rule model.ClassAttendanceRule `json:"rule"`
}

func (v *APIServerV1) coordinatingClassRulesPost(r *http.Request, classId int64) apiResponse {
	var req coordinatingClassRulesPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	ruleString, env, err := req.Verify()
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("rule failed validation: %s", err))
	}

	rule, err := v.db.CreateNewCoordinatingClassRule(r.Context(), database.CreateNewCoordinatingClassRuleParams{
		ClassID:     classId,
		CreatorID:   oauth2.GetAuthContext(r.Context()).User.ID,
		Title:       req.Title,
		Description: req.Description,
		Rule:        ruleString,
		Env:         env,
	})
	if err != nil {
		switch {
		case errors.Is(err, qrm.ErrNoRows):
			return newErrorResponse(http.StatusBadRequest, "not allowed to create new rule")
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusConflict, "rule with same title already exists")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process coordinating class post database action")
		}
	}

	return coordinatingClassRulesPostResponse{
		newSuccessResponse(),
		rule,
	}
}
