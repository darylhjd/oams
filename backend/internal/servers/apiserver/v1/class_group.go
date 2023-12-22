package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) classGroup(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	groupId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, classGroupUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupGet(r, groupId)
	case http.MethodPatch:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGroupGetResponse struct {
	response
	ClassGroup model.ClassGroup `json:"class_group"`
}

func (v *APIServerV1) classGroupGet(r *http.Request, id int64) apiResponse {
	group, err := v.db.GetClassGroup(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class group does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group get database action")
	}

	return classGroupGetResponse{
		newSuccessResponse(),
		group,
	}
}
