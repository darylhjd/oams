package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/middleware"
)

func Test_user(t *testing.T) {
	tests := []struct {
		name        string
		withAccount any
		wantCode    int
		wantBody    string
	}{
		{
			"request with account in context",
			confidential.Account{HomeAccountID: uuid.NewString(), PreferredUsername: uuid.NewString()},
			http.StatusOK,
			"",
		},
		{
			"request with wrong account type in context",
			time.Time{},
			http.StatusInternalServerError,
			"unexpected account data type",
		},
		{
			"request with no account in context",
			nil,
			http.StatusInternalServerError,
			"unexpected account data type",
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := newTestAPIServerV1(t)

			req := httptest.NewRequest(http.MethodGet, userUrl, nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.AccountContextKey, tt.withAccount))
			rr := httptest.NewRecorder()
			v1.user(rr, req)

			a.Equal(tt.wantCode, rr.Code)
			a.Contains(string(rr.Body.Bytes()), tt.wantBody)
			if tt.wantCode != http.StatusOK {
				return
			}

			expectedBody, err := json.Marshal(userResponse{
				HomeAccountID:     tt.withAccount.(confidential.Account).HomeAccountID,
				PreferredUsername: tt.withAccount.(confidential.Account).PreferredUsername,
			})
			a.Nil(err)
			a.Contains(string(rr.Body.Bytes()), string(expectedBody))
		})
	}
}
