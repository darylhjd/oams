package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/jackc/pgx/v5"
)

func (v *APIServerV1) classGroupSession(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	sessionId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, classGroupSessionUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, classGroupSessionUrl, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group session id"))
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

	v.writeResponse(w, classGroupSessionUrl, resp)
}

type classGroupSessionGetResponse struct {
	response
	ClassGroupSession database.ClassGroupSession `json:"class_group_session"`
}

func (v *APIServerV1) classGroupSessionGet(r *http.Request, id int64) apiResponse {
	session, err := v.db.Q.GetClassGroupSession(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class group session does not exist")
		}
	}

	return classGroupSessionGetResponse{
		newSuccessResponse(),
		session,
	}
}
