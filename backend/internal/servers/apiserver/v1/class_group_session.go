package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) classGroupSession(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	sessionId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, classGroupSessionUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group session id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupSessionGet(r, sessionId)
	case http.MethodPatch:
		resp = v.classGroupSessionPatch(r, sessionId)
	case http.MethodDelete:
		resp = v.classGroupSessionDelete(r, sessionId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGroupSessionGetResponse struct {
	response
	ClassGroupSession model.ClassGroupSession `json:"class_group_session"`
}

func (v *APIServerV1) classGroupSessionGet(r *http.Request, id int64) apiResponse {
	session, err := v.db.GetClassGroupSession(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class group session does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group session get database action")
	}

	return classGroupSessionGetResponse{
		newSuccessResponse(),
		session,
	}
}

type classGroupSessionPatchRequest struct {
	ClassGroupSession database.UpdateClassGroupSessionParams `json:"class_group_session"`
}

type classGroupSessionPatchResponse struct {
	response
	ClassGroupSession classGroupSessionPatchClassGroupSessionResponseFields `json:"class_group_session"`
}

type classGroupSessionPatchClassGroupSessionResponseFields struct {
	ID           int64     `json:"id"`
	ClassGroupID int64     `json:"class_group_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Venue        string    `json:"venue"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (v *APIServerV1) classGroupSessionPatch(r *http.Request, id int64) apiResponse {
	var req classGroupSessionPatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	session, err := v.db.UpdateClassGroupSession(r.Context(), id, req.ClassGroupSession)
	if err != nil {
		switch {
		case errors.Is(err, qrm.ErrNoRows):
			return newErrorResponse(http.StatusNotFound, "class group session to update does not exist")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusBadRequest, "class_group_id does not exist")
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusConflict, "class group session with same class_group_id and start_time already exists")
		case database.ErrSQLState(err, database.SQLStateCheckConstraintFailure):
			return newErrorResponse(http.StatusBadRequest, "class group session cannot have a start_time later than or equal to end_time")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process class group session patch database action")
		}
	}

	return classGroupSessionPatchResponse{
		newSuccessResponse(),
		classGroupSessionPatchClassGroupSessionResponseFields{
			session.ID,
			session.ClassGroupID,
			session.StartTime,
			session.EndTime,
			session.Venue,
			session.UpdatedAt,
		},
	}
}

type classGroupSessionDeleteResponse struct {
	response
}

func (v *APIServerV1) classGroupSessionDelete(r *http.Request, id int64) apiResponse {
	_, err := v.db.DeleteClassGroupSession(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, qrm.ErrNoRows):
			return newErrorResponse(http.StatusNotFound, "class group session to delete does not exist")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusConflict, "class group session to delete is still referenced")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process class group session delete database action")
		}
	}

	return classGroupSessionDeleteResponse{newSuccessResponse()}
}
