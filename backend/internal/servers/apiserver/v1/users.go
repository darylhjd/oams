package v1

import (
	"fmt"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) users(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.usersGet(r)
	case http.MethodPost:
		resp = v.usersPost(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type usersGetResponse struct {
	response
	Users []database.User `json:"users"`
}

func (v *APIServerV1) usersGet(r *http.Request) apiResponse {
	users, err := v.db.Q.ListUsers(r.Context())
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process users get database action")
	}

	resp := usersGetResponse{
		newSuccessResponse(),
		make([]database.User, 0, len(users)),
	}

	resp.Users = append(resp.Users, users...)
	return resp
}

type usersPostRequest struct {
	User database.CreateUserParams `json:"user"`
}

type usersPostResponse struct {
	response
	User database.CreateUserRow `json:"user"`
}

func (v *APIServerV1) usersPost(r *http.Request) apiResponse {
	var req usersPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	if req.User.ID == sessionUserId {
		return newErrorResponse(http.StatusUnprocessableEntity, "id is not allowed")
	}

	user, err := v.db.Q.CreateUser(r.Context(), req.User)
	if err != nil {
		if database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex) {
			return newErrorResponse(http.StatusConflict, "user with same id already exists")

		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process users post database action")
	}

	return usersPostResponse{
		newSuccessResponse(),
		user,
	}
}
