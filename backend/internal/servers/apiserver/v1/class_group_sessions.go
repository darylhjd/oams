package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
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
		v.logInternalServerError(r, err)
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
	ClassGroupSession classGroupSessionsPostClassGroupSessionRequestFields `json:"class_group_session"`
}

type classGroupSessionsPostClassGroupSessionRequestFields struct {
	ClassGroupID int64  `json:"class_group_id"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	Venue        string `json:"venue"`
}

func (r classGroupSessionsPostRequest) createClassGroupSessionParams() database.CreateClassGroupSessionParams {
	return database.CreateClassGroupSessionParams{
		ClassGroupID: r.ClassGroupSession.ClassGroupID,
		StartTime:    pgtype.Timestamptz{Time: time.UnixMicro(r.ClassGroupSession.StartTime), Valid: true},
		EndTime:      pgtype.Timestamptz{Time: time.UnixMicro(r.ClassGroupSession.EndTime), Valid: true},
		Venue:        r.ClassGroupSession.Venue,
	}
}

type classGroupSessionsPostResponse struct {
	response
	ClassGroupSession database.CreateClassGroupSessionRow `json:"class_group_session"`
}

func (v *APIServerV1) classGroupSessionsPost(r *http.Request) apiResponse {
	var req classGroupSessionsPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	session, err := v.db.Q.CreateClassGroupSession(r.Context(), req.createClassGroupSessionParams())
	if err != nil {
		switch {
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusBadRequest, "class_group_id does not exist")
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusConflict, "class group session with same class_group_id and start_time already exists")
		case database.ErrSQLState(err, database.SQLStateCheckConstraintFailure):
			return newErrorResponse(http.StatusBadRequest, "class group session cannot have a start_time later than or equal to end_time")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process class group sessions post database action")
		}
	}

	return classGroupSessionsPostResponse{
		newSuccessResponse(),
		session,
	}
}
