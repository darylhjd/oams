package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/middleware"
)

func TestAPIServerV1_users(t *testing.T) {
	tests := []struct {
		name            string
		withAuthContext any
		wantCode        int
	}{
		{
			"request with account in context",
			middleware.AuthContext{
				AuthResult: oauth2.NewMockAzureAuthenticator().MockAuthResult(),
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
			http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			v1 := newTestAPIServerV1(t)

			req := httptest.NewRequest(http.MethodGet, usersUrl, nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.AuthContextKey, tt.withAuthContext))
			rr := httptest.NewRecorder()
			v1.users(rr, req)

			a.Equal(tt.wantCode, rr.Code)
			if tt.wantCode != http.StatusOK {
				return
			}

			var actualResp usersResponse
			err := json.Unmarshal(rr.Body.Bytes(), &actualResp)
			a.Nil(err)

			var session *database.Student
			if tt.withAuthContext != nil {
				authResult := tt.withAuthContext.(middleware.AuthContext).AuthResult
				session = &database.Student{
					ID: authResult.IDToken.Name,
					Email: pgtype.Text{
						String: authResult.Account.PreferredUsername,
						Valid:  true,
					},
					CreatedAt: actualResp.SessionUser.CreatedAt,
					UpdatedAt: actualResp.SessionUser.UpdatedAt,
				}
			}

			expectedResp := usersResponse{
				SessionUser: session,
				Users:       []database.Student{},
			}
			a.Equal(expectedResp, actualResp)
		})
	}
}
