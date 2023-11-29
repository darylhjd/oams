package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
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
	Users []model.User `json:"users"`
}

func (v *APIServerV1) usersGet(r *http.Request) apiResponse {
	params, err := database.DecodeListQueryParams(
		r.URL.Query(), table.Users, table.Users.AllColumns,
	)
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}

	users, err := v.db.ListUsers(r.Context(), params)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process users get database action")
	}

	resp := usersGetResponse{
		newSuccessResponse(),
		make([]model.User, 0, len(users)),
	}

	resp.Users = append(resp.Users, users...)
	return resp
}

type usersPostRequest struct {
	User database.CreateUserParams `json:"user"`
}

type usersPostResponse struct {
	response
	User usersPostUserResponseFields `json:"user"`
}

type usersPostUserResponseFields struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Role      model.UserRole `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
}

func (v *APIServerV1) usersPost(r *http.Request) apiResponse {
	var req usersPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	if req.User.ID == sessionUserId {
		return newErrorResponse(http.StatusUnprocessableEntity, "id is not allowed")
	}

	user, err := v.db.CreateUser(r.Context(), req.User)
	if err != nil {
		if database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex) {
			return newErrorResponse(http.StatusConflict, "user with same id already exists")

		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process users post database action")
	}

	return usersPostResponse{
		newSuccessResponse(),
		usersPostUserResponseFields{
			user.ID,
			user.Name,
			user.Email,
			user.Role,
			user.CreatedAt,
		},
	}
}
