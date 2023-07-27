package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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
	Users []usersGetUserResponseFields `json:"users"`
}

type usersGetUserResponseFields struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Email     *string           `json:"email"`
	Role      database.UserRole `json:"role"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func (r usersGetResponse) fromDatabaseUsers(users []database.User) usersGetResponse {
	resp := usersGetResponse{
		response: newSuccessResponse(),
		Users:    make([]usersGetUserResponseFields, 0, len(users)),
	}

	for _, user := range users {
		respUser := usersGetUserResponseFields{
			ID:        user.ID,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
		}

		if user.Email.Valid {
			respUser.Email = &user.Email.String
		}

		resp.Users = append(resp.Users, respUser)
	}

	return resp
}

func (v *APIServerV1) usersGet(r *http.Request) apiResponse {
	users, err := v.db.Q.ListUsers(r.Context())
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, "could not process users get database action")
	}

	return usersGetResponse{}.fromDatabaseUsers(users)
}

type usersCreateRequest struct {
	User usersCreateUserRequestFields `json:"user"`
}

type usersCreateUserRequestFields struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	Email *string           `json:"email"`
	Role  database.UserRole `json:"role"`
}

func (r usersCreateRequest) createUserParams() database.CreateUserParams {
	params := database.CreateUserParams{
		ID:   r.User.ID,
		Name: r.User.Name,
		Role: r.User.Role,
	}

	if r.User.Email != nil {
		params.Email = pgtype.Text{String: *r.User.Email, Valid: true}
	}

	return params
}

type usersCreateResponse struct {
	response
	User usersCreateUserResponseFields `json:"user"`
}

type usersCreateUserResponseFields struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Email     *string           `json:"email"`
	Role      database.UserRole `json:"role"`
	CreatedAt time.Time         `json:"created_at"`
}

func (r usersCreateResponse) fromDatabaseCreateUserRow(row database.CreateUserRow) usersCreateResponse {
	resp := usersCreateResponse{
		newSuccessResponse(),
		usersCreateUserResponseFields{
			ID:        row.ID,
			Name:      row.Name,
			Role:      row.Role,
			CreatedAt: row.CreatedAt.Time,
		},
	}

	if row.Email.Valid {
		resp.User.Email = &row.Email.String
	}

	return resp
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

	user, err := v.db.Q.CreateUser(r.Context(), req.createUserParams())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == database.SQLStateDuplicateKeyOrIndex {
			return newErrorResponse(http.StatusConflict, "user with same id already exists")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process users post database action")
	}

	return usersCreateResponse{}.fromDatabaseCreateUserRow(user)
}
