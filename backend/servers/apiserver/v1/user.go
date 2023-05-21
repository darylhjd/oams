package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/internal/middleware"
)

type userResponse struct {
	HomeAccountID     string `json:"home_account_id"`
	PreferredUsername string `json:"username"`
}

func (v *APIServerV1) user(w http.ResponseWriter, r *http.Request) {
	v.l.Debug(fmt.Sprintf("%s - received user request", namespace))

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
