package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
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
				[]model.User{
					{
						ID:   "EXISTING_USER",
						Role: model.UserRole_User,
					},
				},
			},
		},
		{
			"request with no user in database",
			false,
			usersGetResponse{
				newSuccessResponse(),
				[]model.User{},
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
					createdUser := tests.StubUser(t, ctx, v1.db, user.ID, user.Role)
					userPtr := &tt.wantResponse.Users[idx]
					userPtr.CreatedAt, userPtr.UpdatedAt = createdUser.CreatedAt, createdUser.CreatedAt
				}
			}

			req := httptest.NewRequest(http.MethodGet, usersUrl, nil)
			actualResp, ok := v1.usersGet(req).(usersGetResponse)
			a.True(ok)
			a.Equal(tt.wantResponse, actualResp)
		})
	}
}

func TestAPIServerV1_usersGetQueryParams(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name           string
		query          url.Values
		wantStatusCode int
		wantErr        string
	}{
		{
			"sort with correct column",
			url.Values{
				"sort": []string{"role"},
			},
			http.StatusOK,
			"",
		},
		{
			"sort with wrong column",
			url.Values{
				"sort": []string{"wrong"},
			},
			http.StatusBadRequest,
			"unknown sort column `wrong`",
		},
		{
			"sort with no value",
			url.Values{
				"sort": []string{},
			},
			http.StatusOK,
			"",
		},
		{
			"limit present",
			url.Values{
				"limit": []string{"1"},
			},
			http.StatusOK,
			"",
		},
		{
			"limit with no value",
			url.Values{
				"limit": []string{},
			},
			http.StatusOK,
			"",
		},
		{
			"offset present",
			url.Values{
				"offset": []string{"1"},
			},
			http.StatusOK,
			"",
		},
		{
			"offset with no value",
			url.Values{
				"offset": []string{},
			},
			http.StatusOK,
			"",
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

			u := url.URL{Path: usersUrl}
			u.RawQuery = tt.query.Encode()

			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			resp := v1.usersGet(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				_, ok := resp.(usersGetResponse)
				a.True(ok)
			}
		})
	}
}

func TestAPIServerV1_usersPost(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name             string
		withRequest      usersPostRequest
		withExistingUser bool
		wantResponse     usersPostResponse
		wantStatusCode   int
		wantErr          string
	}{
		{
			"request with no existing user",
			usersPostRequest{
				database.CreateUserParams{
					ID:   "NEW_USER",
					Role: model.UserRole_User,
				},
			},
			false,
			usersPostResponse{
				newSuccessResponse(),
				usersPostUserResponseFields{
					ID:   "NEW_USER",
					Role: model.UserRole_User,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with existing user",
			usersPostRequest{
				database.CreateUserParams{
					ID:   "EXISTING_USER",
					Role: model.UserRole_User,
				},
			},
			true,
			usersPostResponse{},
			http.StatusConflict,
			"user with same id already exists",
		},
		{
			"request with special session user id me",
			usersPostRequest{
				database.CreateUserParams{
					ID:   sessionUserId,
					Role: model.UserRole_User,
				},
			},
			false,
			usersPostResponse{},
			http.StatusUnprocessableEntity,
			"id is not allowed",
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
				_ = tests.StubUser(t, ctx, v1.db, tt.withRequest.User.ID, tt.withRequest.User.Role)
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPost, usersUrl, bytes.NewReader(reqBodyBytes))
			resp := v1.usersPost(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(usersPostResponse)
				a.True(ok)

				tt.wantResponse.User.CreatedAt = actualResp.User.CreatedAt
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
