package v1

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) classGroupSessions(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupSessionsGet(r)
	case http.MethodPost:
		resp = v.classGroupSessionsPost(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGroupSessionsGetResponse struct {
	response
	ClassGroupSessions []database.ClassGroupSession `json:"class_group_sessions"`
}

func (v *APIServerV1) classGroupSessionsGet(r *http.Request) apiResponse {
	sessions, err := v.db.Q.ListClassGroupSessions(r.Context())
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, "could not process class group sessions get database action")
	}

	resp := classGroupSessionsGetResponse{
		newSuccessResponse(),
		make([]database.ClassGroupSession, 0, len(sessions)),
	}

	resp.ClassGroupSessions = append(resp.ClassGroupSessions, sessions...)
	return resp
}

type classGroupSessionsPostRequest struct {
	ClassGroupSession database.CreateClassGroupSessionParams `json:"class_group_session"`
}

type classGroupSessionsPostResponse struct {
	response
	ClassGroupSession database.CreateClassGroupSessionRow `json:"class_group_session"`
}

func (v *APIServerV1) classGroupSessionsPost(r *http.Request) apiResponse {
	var (
		b   bytes.Buffer
		req classGroupSessionsPostRequest
	)

	if _, err := b.ReadFrom(r.Body); err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if err := json.Unmarshal(b.Bytes(), &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, "could not parse request body")
	}

	session, err := v.db.Q.CreateClassGroupSession(r.Context(), req.ClassGroupSession)
	if err != nil {
		switch {
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusConflict, "class group session with same class_group_id and start_time already exists")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusBadRequest, "class_group_id is not valid")
		default:
			return newErrorResponse(http.StatusInternalServerError, "could not process class group sessions post database action")
		}
	}

	return classGroupSessionsPostResponse{
		newSuccessResponse(),
		session,
	}
}
