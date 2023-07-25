package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func Test_usersGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name            string
		withAuthContext any
		wantResponse    usersGetResponse
		wantErr         string
	}{
		{
			"request with account in context",
			tests.NewMockAuthContext(),
			usersGetResponse{
				response: newSuccessfulResponse(),
				SessionUser: &database.User{
					ID: tests.MockAuthenticatorIDTokenName,
					Email: pgtype.Text{
						String: tests.MockAuthenticatorAccountPreferredUsername,
						Valid:  true,
					},
				},
				Users: []database.User{},
			},
			"",
		},
		{
			"request with no account in context",
			nil,
			usersGetResponse{
				response: newSuccessfulResponse(),
				Users:    []database.User{},
			},
			"",
		},
		{
			"request with wrong auth context type",
			time.Time{},
			usersGetResponse{},
			middleware.ErrUnexpectedAuthContextType.Error(),
		},
	}

	for _, tt := range tts {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := assert.New(t)
			ctx := context.Background()
			id := uuid.NewString()

			v1 := newTestAPIServerV1(t, id)
			defer tests.TearDown(t, v1.db, id)

			tests.StubAuthContextUser(t, ctx, v1.db.Q)

			req := httptest.NewRequest(http.MethodGet, usersUrl, nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.AuthContextKey, tt.withAuthContext))
			actualResp, err := v1.usersGet(req)

			switch {
			case tt.wantErr != "":
				a.NotNil(err)
				a.Contains(err.Error(), tt.wantErr)
				return
			case tt.withAuthContext != nil:
				a.Equal(tt.wantResponse.SessionUser.ID, actualResp.SessionUser.ID)
				tt.wantResponse.SessionUser = actualResp.SessionUser
			}

			a.Equal(tt.wantResponse, actualResp)
		})
	}
}
