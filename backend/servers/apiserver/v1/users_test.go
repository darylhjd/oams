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
		name             string
		withExistingUser bool
		wantResponse     usersGetResponse
	}{
		{
			"request with user in database",
			true,
			usersGetResponse{
				newSuccessResponse(),
				[]database.User{
					{
						ID:   "EXISTING_USER",
						Role: database.UserRoleSTUDENT,
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

			if tt.withExistingUser {
				for idx, user := range tt.wantResponse.Users {
					createdUser := tests.StubUser(t, ctx, v1.db.Q, user.ID, user.Role)
					userPtr := &tt.wantResponse.Users[idx]
					userPtr.CreatedAt, userPtr.UpdatedAt = createdUser.CreatedAt, createdUser.CreatedAt
				}
			}

			req := httptest.NewRequest(http.MethodGet, usersUrl, nil)
			actualResp, ok := v1.usersGet(req).(usersGetResponse)
			a.True(ok)

			a.Equal(len(tt.wantResponse.Users), len(actualResp.Users))
			a.Equal(tt.wantResponse, actualResp)
		})
	}
}

func TestAPIServerV1_usersPost(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name             string
		withRequest      usersCreateRequest
		withExistingUser bool
		wantResponse     usersCreateResponse
		wantStatusCode   int
		wantErr          string
	}{
		{
			"request with no existing user",
			usersCreateRequest{
				database.CreateUserParams{
					ID:   "NEW_USER",
					Role: database.UserRoleSTUDENT,
				},
			},
			false,
			usersCreateResponse{
				newSuccessResponse(),
				database.CreateUserRow{
					ID:   "NEW_USER",
					Role: database.UserRoleSTUDENT,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with existing user",
			usersCreateRequest{
				database.CreateUserParams{
					ID:   "EXISTING_USER",
					Role: database.UserRoleSTUDENT,
				},
			},
			true,
			usersCreateResponse{},
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

			if tt.withExistingUser {
				_ = tests.StubUser(t, ctx, v1.db.Q, tt.withRequest.User.ID, tt.withRequest.User.Role)
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPost, usersUrl, bytes.NewReader(reqBodyBytes))
			resp := v1.usersCreate(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(usersCreateResponse)
				a.True(ok)

				tt.wantResponse.User.CreatedAt = actualResp.User.CreatedAt
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
