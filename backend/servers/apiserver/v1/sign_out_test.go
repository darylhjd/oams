package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

func Test_signOut(t *testing.T) {
	tests := []struct {
		name            string
		withAuthContext any
		wantCode        int
	}{
		{
			"request with account in context",
			middleware.AuthContext{
				AuthResult: confidential.AuthResult{
					Account: confidential.Account{HomeAccountID: uuid.NewString(), PreferredUsername: uuid.NewString()},
				},
			},
			http.StatusOK,
		},
		{
			"request with wrong account type in context",
			time.Time{},
			http.StatusInternalServerError,
		},
		{
			"request with no account in context",
			nil,
			http.StatusInternalServerError,
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := newTestAPIServerV1(t)

			req := httptest.NewRequest(http.MethodGet, signOutUrl, nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.AuthContextKey, tt.withAuthContext))
			rr := httptest.NewRecorder()
			v1.signOut(rr, req)

			a.Equal(tt.wantCode, rr.Code)
			if tt.wantCode != http.StatusOK {
				return
			}

			for _, cookie := range rr.Result().Cookies() {
				if cookie.Name == oauth2.SessionCookieIdent {
					return
				}
			}
			a.FailNow("could not detect expected session cookie")
		})
	}
}
