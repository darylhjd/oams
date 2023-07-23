package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/middleware"
)

// usersResponse is a struct detailing the response body of the users endpoint.
type usersResponse struct {
	SessionUser *database.Student  `json:"session_user"`
	Users       []database.Student `json:"users"`
}

// users endpoint returns useful information on the current session user and information on any requested users..
func (v *APIServerV1) users(w http.ResponseWriter, r *http.Request) {
	resp := usersResponse{
		Users: []database.Student{},
	}
	// Get users data.
	authContext, isSignedIn, err := middleware.GetAuthContext(r)
	switch {
	case err != nil:
		http.Error(w, "error getting current user auth session", http.StatusInternalServerError)
		return
	case isSignedIn:
		student, err := v.db.Q.GetStudent(r.Context(), authContext.AuthResult.IDToken.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp.SessionUser = &student
	}

	v.l.Debug("response", zap.Any("response", resp))
	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(bytes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("url", usersUrl),
			zap.Error(err))
	}
}
