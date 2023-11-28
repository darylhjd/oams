package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) classManager(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	managerId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, classManagerUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class manager id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classManagerGet(r, managerId)
	case http.MethodPatch:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classManagerGetResponse struct {
	response
	ClassManager model.ClassManager `json:"class_manager"`
}

func (v *APIServerV1) classManagerGet(r *http.Request, id int64) apiResponse {
	manager, err := v.db.GetClassManager(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class manager does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class manager get database action")
	}

	return classManagerGetResponse{
		newSuccessResponse(),
		manager,
	}
}
