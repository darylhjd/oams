package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) classGroupManager(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	managerId, err := strconv.ParseInt(r.PathValue("managerId"), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group manager id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupManagerGet(r, managerId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGroupManagerGetResponse struct {
	response
	Manager model.ClassGroupManager `json:"manager"`
}

func (v *APIServerV1) classGroupManagerGet(r *http.Request, managerId int64) apiResponse {
	manager, err := v.db.GetClassGroupManager(r.Context(), managerId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class group manager does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group manager get database action")
	}

	return classGroupManagerGetResponse{
		newSuccessResponse(),
		manager,
	}
}
