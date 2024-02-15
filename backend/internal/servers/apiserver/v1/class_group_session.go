package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) classGroupSession(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	sessionId, err := strconv.ParseInt(r.PathValue("sessionId"), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group session id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupSessionGet(r, sessionId)
	case http.MethodPatch:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGroupSessionGetResponse struct {
	response
	ClassGroupSession model.ClassGroupSession `json:"class_group_session"`
}

func (v *APIServerV1) classGroupSessionGet(r *http.Request, sessionId int64) apiResponse {
	s, err := v.db.GetClassGroupSession(r.Context(), sessionId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class group session does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group session get database action")
	}

	return classGroupSessionGetResponse{
		newSuccessResponse(),
		s,
	}
}
