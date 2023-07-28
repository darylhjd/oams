package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/jackc/pgx/v5"
)

func (v *APIServerV1) class(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, classUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, classUrl, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classGet(r, classId)
	case http.MethodPut:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, classUrl, resp)
}

type classGetResponse struct {
	response
	Class database.Class `json:"class"`
}

func (v *APIServerV1) classGet(r *http.Request, id int64) apiResponse {
	class, err := v.db.Q.GetClass(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process class get database action")
	}

	return classGetResponse{
		newSuccessResponse(),
		class,
	}
}
