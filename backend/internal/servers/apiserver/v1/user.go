package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	"github.com/go-jet/jet/v2/qrm"
)

const (
	sessionUserId = "me"
)

func (v *APIServerV1) user(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	userId := strings.TrimPrefix(r.URL.Path, userUrl)
	switch m := r.Method; {
	case userId == sessionUserId:
		resp = v.userMe(r)
	case m == http.MethodGet:
		resp = v.userGet(r, userId)
	case m == http.MethodPatch:
		resp = v.userPatch(r, userId)
	case m == http.MethodDelete:
		resp = v.userDelete(r, userId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type userMeResponse struct {
	response
	SessionUser                model.User                           `json:"session_user"`
	UpcomingClassGroupSessions []database.UpcomingClassGroupSession `json:"upcoming_class_group_sessions"`
}

func (v *APIServerV1) userMe(r *http.Request) apiResponse {
	resp := userMeResponse{
		response:                   newSuccessResponse(),
		UpcomingClassGroupSessions: []database.UpcomingClassGroupSession{},
	}

	authContext := oauth2.GetAuthContext(r.Context())

	sessionUser, err := v.db.GetUser(r.Context(), authContext.User.ID)
	if err != nil {
		v.logInternalServerError(r, fmt.Errorf("expected session user in database: %w", err))
		return newErrorResponse(http.StatusInternalServerError, "could get session user from database")
	}

	upcomingClassGroupSessions, err := v.db.GetUserUpcomingClassGroupSessions(r.Context(), sessionUser.ID)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get session user upcoming class group sessions")
	}

	resp.SessionUser = sessionUser
	resp.UpcomingClassGroupSessions = append(resp.UpcomingClassGroupSessions, upcomingClassGroupSessions...)
	return resp
}

type userGetResponse struct {
	response
	User model.User `json:"user"`
}

func (v *APIServerV1) userGet(r *http.Request, id string) apiResponse {
	user, err := v.db.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested user does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process user get database action")
	}

	return userGetResponse{
		newSuccessResponse(),
		user,
	}
}

type userPatchRequest struct {
	User database.UpdateUserParams `json:"user"`
}

type userPatchResponse struct {
	response
	User userPatchUserResponseFields `json:"user"`
}

type userPatchUserResponseFields struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Role      model.UserRole `json:"role"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (v *APIServerV1) userPatch(r *http.Request, id string) apiResponse {
	var req userPatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	user, err := v.db.UpdateUser(r.Context(), id, req.User)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "user to update does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process user patch database action")
	}

	return userPatchResponse{
		newSuccessResponse(),
		userPatchUserResponseFields{
			user.ID,
			user.Name,
			user.Email,
			user.Role,
			user.UpdatedAt,
		},
	}
}

type userDeleteResponse struct {
	response
}

func (v *APIServerV1) userDelete(r *http.Request, id string) apiResponse {
	_, err := v.db.DeleteUser(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, qrm.ErrNoRows):
			return newErrorResponse(http.StatusNotFound, "user to delete does not exist")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusConflict, "user to delete is still referenced")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process user delete database action")
		}
	}

	return userDeleteResponse{newSuccessResponse()}
}
