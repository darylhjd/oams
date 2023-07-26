package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/jackc/pgx/v5"
)

// users endpoint returns useful information on the current session user and information on any requested users.
func (v *APIServerV1) user(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	userId := strings.TrimPrefix(r.URL.Path, userUrl)
	switch r.Method {
	case http.MethodGet:
		resp = v.userGet(r, userId)
	case http.MethodPut:
		resp = v.userPut(r)
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, fmt.Sprintf("method %s is not allowed", r.Method))
	}

	v.writeResponse(w, userUrl, resp)
}

type userGetResponse struct {
	response
	User database.User `json:"user"`
}

func (v *APIServerV1) userGet(r *http.Request, id string) apiResponse {
	resp := userGetResponse{
		response: newSuccessResponse(),
	}

	user, err := v.db.Q.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested user does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	resp.User = user
	return resp
}

type userPutRequest struct {
	Users []database.UpsertUsersParams `json:"users"`
}

type userPutResponse struct {
	response
	Users []database.User `json:"users"`
}

func (v *APIServerV1) userPut(r *http.Request) apiResponse {
	var (
		b   bytes.Buffer
		req userPutRequest
	)
	resp := userPutResponse{
		newSuccessResponse(),
		[]database.User{},
	}

	if _, err := b.ReadFrom(r.Body); err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if err := json.Unmarshal(b.Bytes(), &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, "could not parse request body")
	}

	users, err := upsertUsers(r.Context(), v.db, nil, req.Users)
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, "could not process users database action")
	}

	resp.Users = users
	return resp
}

// upsertUsers inserts the provided usersParams into the specified database. If tx is nil, a new transaction is started.
// Otherwise, a nested transaction (using save points) is used.
func upsertUsers(ctx context.Context, db *database.DB, tx pgx.Tx, usersParams []database.UpsertUsersParams) ([]database.User, error) {
	var err error

	if tx != nil {
		tx, err = tx.Begin(ctx)
	} else {
		tx, err = db.C.Begin(ctx)
	}

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	q := db.Q.WithTx(tx)

	if err = q.UpsertUsers(ctx, usersParams).Close(); err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(usersParams))
	for _, param := range usersParams {
		ids = append(ids, param.ID)
	}

	users, err := q.GetUsersByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return users, tx.Commit(ctx)
}
