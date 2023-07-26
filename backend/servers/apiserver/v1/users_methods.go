package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/jackc/pgx/v5"
)

type usersGetResponse struct {
	response
	SessionUser *database.User  `json:"session_user"`
	Users       []database.User `json:"users"`
}

type usersGetQueries struct {
	ids []string
}

const (
	usersGetQueriesIdKey = "ids"
)

func (v *APIServerV1) usersGet(r *http.Request) apiResponse {
	resp := usersGetResponse{
		response: newSuccessResponse(),
		Users:    []database.User{},
	}

	// Fill session user.
	authContext, isSignedIn, err := middleware.GetAuthContext(r)
	switch {
	case err != nil:
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	case isSignedIn:
		student, err := v.db.Q.GetUser(r.Context(), authContext.AuthResult.IDToken.Name)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, err.Error())
		}

		resp.SessionUser = &student
	}

	// Parse queries
	var queries usersGetQueries
	{
		q := r.URL.Query()
		queries.ids = q[usersGetQueriesIdKey]
	}

	students, err := v.db.Q.GetUsersByIDs(r.Context(), queries.ids)
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	resp.Users = append(resp.Users, students...)
	return resp
}

type usersPutRequest struct {
	Users []database.UpsertUsersParams `json:"users"`
}

type usersPutResponse struct {
	response
	Users []database.User `json:"users"`
}

func (v *APIServerV1) usersPut(r *http.Request) apiResponse {
	var (
		b   bytes.Buffer
		req usersPutRequest
	)
	resp := usersPutResponse{
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
