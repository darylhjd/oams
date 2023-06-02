package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/middleware"
)

// userResponse is a struct detailing the response body of the user endpoint.
type userResponse struct {
	HomeAccountID     string `json:"home_account_id"`
	PreferredUsername string `json:"username"`
}

// user endpoint returns useful information on a User.
func (v *APIServerV1) user(w http.ResponseWriter, r *http.Request) {
	// Get user data.
	acct, ok := r.Context().Value(middleware.AccountContextKey).(confidential.Account)
	if !ok {
		http.Error(w, "unexpected account data type", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(&userResponse{
		HomeAccountID:     acct.HomeAccountID,
		PreferredUsername: acct.PreferredUsername,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(bytes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("url", userUrl),
			zap.Error(err))
	}
}
