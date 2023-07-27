package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/middleware"
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
	case m == http.MethodPut:
		resp = v.userPut(r, userId)
	case m == http.MethodDelete:
		resp = v.userDelete(r, userId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, userUrl, resp)
}

type userMeResponse struct {
	response
	SessionUser database.User `json:"session_user"`
}

func (v *APIServerV1) userMe(r *http.Request) apiResponse {
	resp := userMeResponse{response: newSuccessResponse()}

	// Fill session user.
	authContext, isSignedIn, err := middleware.GetAuthContext(r)
	switch {
	case err != nil:
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	case isSignedIn:
		user, err := v.db.Q.GetUser(r.Context(), authContext.AuthResult.IDToken.Name)
		if err != nil {
			return newErrorResponse(http.StatusInternalServerError, "could get session user from database")
		}

		resp.SessionUser = user
	default:
		return newErrorResponse(http.StatusUnauthorized, "client lacks authentication credentials")
	}

	return resp
}

type userGetResponse struct {
	response
	User database.User `json:"user"`
}

func (v *APIServerV1) userGet(r *http.Request, id string) apiResponse {
	user, err := v.db.Q.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested user does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process user get database action")
	}

	return userGetResponse{
		newSuccessResponse(),
		user,
	}
}

type userPutRequest struct {
	User userPutUserRequestFields `json:"user"`
}

type userPutUserRequestFields struct {
	Name  *string            `json:"name"`
	Email *string            `json:"email"`
	Role  *database.UserRole `json:"role"`
}

func (r userPutRequest) updateUserParams(userId string) database.UpdateUserParams {
	params := database.UpdateUserParams{ID: userId}

	// Parse request.
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

type userPutResponse struct {
	response
	User database.UpdateUserRow `json:"user"`
}

func (v *APIServerV1) userPut(r *http.Request, id string) apiResponse {
	var (
		b   bytes.Buffer
		req userPutRequest
	)

	if _, err := b.ReadFrom(r.Body); err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if err := json.Unmarshal(b.Bytes(), &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, "could not parse request body")
	}

	user, err := v.db.Q.UpdateUser(r.Context(), req.updateUserParams(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "user to update does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process user put database action")
	}

	return userPutResponse{
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
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "user to delete does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process user delete database action")
	}

	return userDeleteResponse{newSuccessResponse()}
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
