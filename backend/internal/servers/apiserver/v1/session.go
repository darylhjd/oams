package v1

import (
	"fmt"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

type sessionResponse struct {
	response
	Session *session `json:"session"`
}

type session struct {
	User              model.User                 `json:"user"`
	ManagementDetails database.ManagementDetails `json:"management_details"`
}

func (v *APIServerV1) session(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		v.writeResponse(w, r, newErrorResponse(http.StatusMethodNotAllowed, ""))
		return
	}

	if _, _, err := middleware.CheckAuthorizationToken(r, v.auth); err != nil {
		v.writeResponse(w, r, sessionResponse{newSuccessResponse(), nil})
		return
	}

	middleware.MustAuth(v.getSession, v.auth, v.db)(w, r)
}

func (v *APIServerV1) getSession(w http.ResponseWriter, r *http.Request) {
	authContext := oauth2.GetAuthContext(r.Context())

	sessionUser, err := v.db.GetUser(r.Context(), authContext.User.ID)
	if err != nil {
		v.logInternalServerError(r, fmt.Errorf("expected session user in database: %w", err))
		v.writeResponse(w, r, newErrorResponse(http.StatusInternalServerError, "could get session user from database"))
		return
	}

	details, err := v.db.GetManagementDetails(r.Context())
	if err != nil {
		v.logInternalServerError(r, fmt.Errorf("could not get user managed class groups: %w", err))
		v.writeResponse(w, r, newErrorResponse(http.StatusInternalServerError, "could not get session user managed class groups"))
		return
	}

	v.writeResponse(w, r, sessionResponse{
		newSuccessResponse(),
		&session{sessionUser, details},
	})
}
