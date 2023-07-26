package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/middleware"
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
