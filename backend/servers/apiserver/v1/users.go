package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/middleware"
)

type usersResponse struct {
	response
	SessionUser *database.Student  `json:"session_user"`
	Users       []database.Student `json:"users"`
}

// users endpoint returns useful information on the current session user and information on any requested users..
func (v *APIServerV1) users(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	// Get users data.
	var session *database.Student

	authContext, isSignedIn, err := middleware.GetAuthContext(r)
	switch {
	case err != nil:
		resp = newErrorResponse(http.StatusInternalServerError, err.Error())
	case isSignedIn:
		student, err := v.db.Q.GetStudent(r.Context(), authContext.AuthResult.IDToken.Name)
		if err != nil {
			resp = newErrorResponse(http.StatusInternalServerError, err.Error())
			break
		}

		session = &student
		fallthrough
	default:
		// TODO: Allow request for other users.
		resp = usersResponse{
			response:    newSuccessfulResponse(),
			SessionUser: session,
			Users:       []database.Student{},
		}
	}

	v.writeResponse(w, usersUrl, resp)
}
