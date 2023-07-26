package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_users(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name           string
		withMethod     string
		wantStatusCode int
	}{
		{
			"with GET method",
			http.MethodGet,
			http.StatusOK,
		},
		{
			"with POST method",
			http.MethodPost,
			http.StatusBadRequest,
		},
		{
			"with DELETE method",
			http.MethodDelete,
			http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tts {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := assert.New(t)
			id := uuid.NewString()

			v1 := newTestAPIServerV1(t, id)
			defer tests.TearDown(t, v1.db, id)

			req := httptest.NewRequest(tt.withMethod, usersUrl, nil)
			rr := httptest.NewRecorder()
			v1.users(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_usersGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                    string
		withStubAuthContextUser bool
		wantResponse            usersGetResponse
	}{
		{
			"request with user in database",
			true,
			usersGetResponse{
				newSuccessResponse(),
				[]database.User{
					{
						ID: tests.MockAuthenticatorIDTokenName,
						Email: pgtype.Text{
							String: tests.MockAuthenticatorAccountPreferredUsername,
							Valid:  true,
						},
					},
				},
			},
		},
		{
			"request with no user in database",
			false,
			usersGetResponse{
				response: newSuccessResponse(),
				Users:    []database.User{},
			},
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

			if tt.withStubAuthContextUser {
				tests.StubAuthContextUser(t, ctx, v1.db.Q)
			}

			req := httptest.NewRequest(http.MethodGet, usersUrl, nil)
			actualResp, ok := v1.usersGet(req).(usersGetResponse)
			a.True(ok)

			a.Equal(len(tt.wantResponse.Users), len(actualResp.Users))
		})
	}
}

func TestAPIServerV1_usersPost(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                    string
		withStubAuthContextUser bool
		wantStatusCode          int
		wantErr                 string
	}{
		{
			"request with no existing user",
			false,
			http.StatusOK,
			"",
		},
		{
			"request with existing user",
			true,
			http.StatusConflict,
			"user with same id already exists",
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

			if tt.withStubAuthContextUser {
				tests.StubAuthContextUser(t, ctx, v1.db.Q)
			}

			reqBodyBytes, err := json.Marshal(usersCreateRequest{
				database.CreateUserParams{
					ID: tests.MockAuthenticatorIDTokenName,
					Email: pgtype.Text{
						String: tests.MockAuthenticatorAccountPreferredUsername,
						Valid:  true,
					},
					Role: database.UserRoleSTUDENT,
				},
			})
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPost, usersUrl, bytes.NewReader(reqBodyBytes))
			resp := v1.usersCreate(req)

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(usersCreateResponse)
				a.True(ok)
				a.Equal(tests.MockAuthenticatorIDTokenName, actualResp.User.ID)
			}
		})
	}
}
