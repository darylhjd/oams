package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

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
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type userMeResponse struct {
	response
	SessionUser           model.User `json:"session_user"`
	HasManagedClassGroups bool       `json:"has_managed_class_groups"`
}

func (v *APIServerV1) userMe(r *http.Request) apiResponse {
	resp := userMeResponse{
		response: newSuccessResponse(),
	}

	authContext := oauth2.GetAuthContext(r.Context())

	sessionUser, err := v.db.GetUser(r.Context(), authContext.User.ID)
	if err != nil {
		v.logInternalServerError(r, fmt.Errorf("expected session user in database: %w", err))
		return newErrorResponse(http.StatusInternalServerError, "could get session user from database")
	}

	resp.SessionUser = sessionUser
	// TODO: Add tests for this field.
	resp.HasManagedClassGroups, err = v.db.HasManagedClassGroups(r.Context())
	if err != nil {
		v.logInternalServerError(r, fmt.Errorf("could not get user managed class groups: %w", err))
		return newErrorResponse(http.StatusInternalServerError, "could not get session user managed class groups")
	}

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
	User struct {
		Signature *string `json:"signature"`
	} `json:"user"`
}

type userPatchResponse struct {
	response
	User model.User `json:"user"`
}

func (v *APIServerV1) userPatch(r *http.Request, id string) apiResponse {
	var req userPatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse user request body: %s", err))
	}

	if req.User.Signature != nil {
		if err := v.db.UpdateUserSignature(r.Context(), id, *req.User.Signature); err != nil {
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not update user signature")
		}
	}

	return userPatchResponse{
		newSuccessResponse(),
		oauth2.GetAuthContext(r.Context()).User,
	}
}
