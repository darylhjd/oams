package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/middleware"
)

type usersGetResponse struct {
	response
	SessionUser *database.Student  `json:"session_user"`
	Users       []database.Student `json:"users"`
}

type usersGetQueries struct {
	ids []string
}

const (
	usersGetQueriesIdKey = "ids"
)

func (v *APIServerV1) usersGet(r *http.Request) (usersGetResponse, error) {
	resp := usersGetResponse{
		response: newSuccessfulResponse(),
		Users:    []database.Student{},
	}

	// Fill session user.
	authContext, isSignedIn, err := middleware.GetAuthContext(r)
	switch {
	case err != nil:
		return resp, err
	case isSignedIn:
		student, err := v.db.Q.GetStudent(r.Context(), authContext.AuthResult.IDToken.Name)
		if err != nil {
			return resp, err
		}

		resp.SessionUser = &student
	}

	// Parse queries
	var queries usersGetQueries
	{
		q := r.URL.Query()
		queries.ids = q[usersGetQueriesIdKey]
	}

	students, err := v.db.Q.GetStudentsByIDs(r.Context(), queries.ids)
	if err != nil {
		return resp, err
	}

	resp.Users = append(resp.Users, students...)
	return resp, err
}
