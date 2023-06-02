package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/oauth2"
)

// signOut endpoint invalidates a user session. This is done by requesting that the browser
// remove the cookie containing the session information.
func (v *APIServerV1) signOut(w http.ResponseWriter, _ *http.Request) {
	_ = oauth2.DeleteSessionCookie(w)
}
