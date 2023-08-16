package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	SessionUser                model.User                                      `json:"session_user"`
	UpcomingClassGroupSessions []database.GetUserUpcomingClassGroupSessionsRow `json:"upcoming_class_group_sessions"`
}

func (v *APIServerV1) userMe(r *http.Request) apiResponse {
	resp := userMeResponse{response: newSuccessResponse()}

	authContext, isSignedIn, err := middleware.GetAuthContext(r)
	switch {
	case err != nil:
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	case isSignedIn:
		resp.SessionUser, err = v.db.GetUser(r.Context(), authContext.AuthResult.IDToken.Name)
		if err != nil {
			v.logInternalServerError(r, fmt.Errorf("expected session user in database: %w", err))
			return newErrorResponse(http.StatusInternalServerError, "could get session user from database")
		}
	default:
		return newErrorResponse(http.StatusUnauthorized, "client lacks authentication credentials")
	}

	upcomingClassGroupSessions, err := v.db.Q.GetUserUpcomingClassGroupSessions(r.Context(), resp.SessionUser.ID)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get session user upcoming class group sessions")
	}

	resp.UpcomingClassGroupSessions = append(
		make([]database.GetUserUpcomingClassGroupSessionsRow, 0, len(upcomingClassGroupSessions)),
		upcomingClassGroupSessions...,
	)
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
	User userPatchUserRequestFields `json:"user"`
}

type userPatchUserRequestFields struct {
	Name  *string            `json:"name"`
	Email *string            `json:"email"`
	Role  *database.UserRole `json:"role"`
}

func (r userPatchRequest) updateUserParams(userId string) database.UpdateUserParams {
	params := database.UpdateUserParams{ID: userId}

	if r.User.Name != nil {
		params.Name = pgtype.Text{String: *r.User.Name, Valid: true}
	}

	if r.User.Email != nil {
		params.Email = pgtype.Text{String: *r.User.Email, Valid: true}
	}

	if r.User.Role != nil {
		params.Role = database.NullUserRole{UserRole: *r.User.Role, Valid: true}
	}

	return params
}

type userPatchResponse struct {
	response
	User database.UpdateUserRow `json:"user"`
}

func (v *APIServerV1) userPatch(r *http.Request, id string) apiResponse {
	var req userPatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	user, err := v.db.Q.UpdateUser(r.Context(), req.updateUserParams(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "user to update does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process user patch database action")
	}

	return userPatchResponse{
		newSuccessResponse(),
		user,
	}
}

type userDeleteResponse struct {
	response
}

func (v *APIServerV1) userDelete(r *http.Request, id string) apiResponse {
	_, err := v.db.Q.DeleteUser(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
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
