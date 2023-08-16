package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
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
	ClassGroupSessions []model.ClassGroupSession `json:"class_group_sessions"`
}

func (v *APIServerV1) classGroupSessionsGet(r *http.Request) apiResponse {
	sessions, err := v.db.ListClassGroupSessions(r.Context())
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group sessions get database action")
	}

	resp := classGroupSessionsGetResponse{
		newSuccessResponse(),
		make([]model.ClassGroupSession, 0, len(sessions)),
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
		StartTime:    time.UnixMicro(r.ClassGroupSession.StartTime),
		EndTime:      time.UnixMicro(r.ClassGroupSession.EndTime),
		Venue:        r.ClassGroupSession.Venue,
	}
}

type classGroupSessionsPostResponse struct {
	response
	ClassGroupSession classGroupSessionsPostClassGroupSessionFields `json:"class_group_session"`
}

type classGroupSessionsPostClassGroupSessionFields struct {
	ID           int64     `json:"id"`
	ClassGroupID int64     `json:"class_group_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Venue        string    `json:"venue"`
	CreatedAt    time.Time `json:"created_at"`
}

func (v *APIServerV1) classGroupSessionsPost(r *http.Request) apiResponse {
	var req classGroupSessionsPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	session, err := v.db.CreateClassGroupSession(r.Context(), req.createClassGroupSessionParams())
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
		classGroupSessionsPostClassGroupSessionFields{
			ClassGroupID: session.ClassGroupID,
			StartTime:    session.StartTime,
			EndTime:      session.EndTime,
			Venue:        session.Venue,
			CreatedAt:    session.CreatedAt,
		},
	}
}
