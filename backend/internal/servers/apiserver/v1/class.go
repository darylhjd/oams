package v1

import (
	"errors"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/pkg/to"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) class(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classId, err := to.Int64(r.PathValue("classId"))
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classGet(r, classId)
	case http.MethodPatch:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGetResponse struct {
	response
	Class model.Class `json:"class"`
}

func (v *APIServerV1) classGet(r *http.Request, classId int64) apiResponse {
	class, err := v.db.GetClass(r.Context(), classId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class get database action")
	}

	return classGetResponse{
		newSuccessResponse(),
		class,
	}
}
