package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) coordinatingClass(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classId, err := strconv.ParseInt(r.PathValue("classId"), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.coordinatingClassGet(r, classId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type coordinatingClassGetResponse struct {
	response
	CoordinatingClass database.CoordinatingClass `json:"coordinating_class"`
}

func (v *APIServerV1) coordinatingClassGet(r *http.Request, classId int64) apiResponse {
	class, err := v.db.GetCoordinatingClass(r.Context(), classId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested coordinating class does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process coordinating class get database action")
	}

	return coordinatingClassGetResponse{
		newSuccessResponse(),
		class,
	}
}
