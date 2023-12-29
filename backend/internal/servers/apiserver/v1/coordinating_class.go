package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
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
}

type coordinatingClassPostResponse struct {
	response
}

func (v *APIServerV1) coordinatingClassPost(r *http.Request, id int64) apiResponse {
	return newErrorResponse(http.StatusNotImplemented, "")
}
