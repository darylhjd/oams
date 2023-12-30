package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/rules"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) coordinatingClass(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, coordinatingClassUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.coordinatingClassGet(r, classId)
	case http.MethodPost:
		resp = v.coordinatingClassPost(r, classId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type coordinatingClassGetResponse struct {
	response
	CoordinatingClass database.CoordinatingClass  `json:"coordinating_class"`
	Rules             []model.ClassAttendanceRule `json:"rules"`
}

// TODO: Implement tests for this endpoint.
func (v *APIServerV1) coordinatingClassGet(r *http.Request, id int64) apiResponse {
	class, err := v.db.GetCoordinatingClass(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested coordinating class does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process coordinating class get database action")
	}

	rules, err := v.db.GetCoordinatingClassRules(r.Context(), id)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get coordinating class rules")
	}

	return coordinatingClassGetResponse{
		newSuccessResponse(),
		class,
		append(make([]model.ClassAttendanceRule, 0, len(rules)), rules...),
	}
}

type coordinatingClassPostRequest struct {
	rules.RuleParams
}

type coordinatingClassPostResponse struct {
	response
	Rule model.ClassAttendanceRule `json:"rule"`
}

func (v *APIServerV1) coordinatingClassPost(r *http.Request, id int64) apiResponse {
	var req coordinatingClassPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	if err := req.Verify(); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("rule failed validation: %s", err))
	}

	rule, err := v.db.CreateNewCoordinatingClassRule(r.Context(), database.CreateNewCoordinatingClassRuleParams{
		ClassID:     id,
		Title:       req.Title,
		Description: req.Description,
		Rule:        req.Rule,
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

	return coordinatingClassPostResponse{
		newSuccessResponse(),
		rule,
	}
}
