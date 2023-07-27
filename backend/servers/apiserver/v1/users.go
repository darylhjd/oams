package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/jackc/pgx/v5/pgconn"
)

func (v *APIServerV1) users(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.usersGet(r)
	case http.MethodPost:
		resp = v.usersCreate(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, usersUrl, resp)
}

type usersGetResponse struct {
	response
	Users []database.User `json:"users"`
}

func (v *APIServerV1) usersGet(r *http.Request) apiResponse {
	users, err := v.db.Q.ListUsers(r.Context())
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, "could not process users get database action")
	}

	resp := usersGetResponse{
		newSuccessResponse(),
		make([]database.User, 0, len(users)),
	}

	resp.Users = append(resp.Users, users...)
	return resp
}

type usersCreateRequest struct {
	User database.CreateUserParams `json:"user"`
}

type usersCreateResponse struct {
	response
	User database.CreateUserRow `json:"user"`
}

func (v *APIServerV1) usersCreate(r *http.Request) apiResponse {
	var (
		b   bytes.Buffer
		req usersCreateRequest
	)

	if _, err := b.ReadFrom(r.Body); err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if err := json.Unmarshal(b.Bytes(), &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, "could not parse request body")
	}

	user, err := v.db.Q.CreateUser(r.Context(), req.User)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == database.SQLStateDuplicateKeyOrIndex {
			return newErrorResponse(http.StatusConflict, "user with same id already exists")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process users post database action")
	}

	return usersCreateResponse{
		newSuccessResponse(),
		user,
	}
}
