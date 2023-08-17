package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
						Role: model.UserRole_Student,
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

	baseRecords := 100

	limitTts := []struct {
		name            string
		limit           string
		expectedRecords int
	}{
		{
			"limit less than total records",
			"99",
			99,
		},
		{
			"limit equal total records",
			"100",
			100,
		},
		{
			"limit more than total records",
			"101",
			100,
		},
		{
			"limit is 0",
			"0",
			database.ListDefaultLimit,
		},
		{
			"limit is negative",
			"-1",
			database.ListDefaultLimit,
		},
	}

	for _, tt := range limitTts {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := assert.New(t)
			ctx := context.Background()
			id := uuid.NewString()

			v1 := newTestAPIServerV1(t, id)
			defer tests.TearDown(t, v1.db, id)

			for i := 0; i < baseRecords; i++ {
				tests.StubUser(t, ctx, v1.db, uuid.NewString(), model.UserRole_Student)
			}

			u := url.URL{Path: usersUrl}
			values := u.Query()
			values.Set("limit", tt.limit)
			u.RawQuery = values.Encode()

			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			resp, ok := v1.usersGet(req).(usersGetResponse)
			a.True(ok)
			a.Equal(tt.expectedRecords, len(resp.Users))
		})
	}

	offsetTts := []struct {
		name        string
		offset      string
		wantUsers   bool
		wantFirstID string
	}{
		{
			"offset less than total records",
			"50",
			true,
			"051",
		},
		{
			"offset one less than total records",
			"99",
			true,
			"100",
		},
		{
			"offset equal total records",
			"100",
			false,
			"",
		},
		{
			"offset more than total records",
			"101",
			false,
			"",
		},
		{
			"offset is 0",
			"0",
			true,
			"001",
		},
		{
			"offset is negative",
			"-1",
			true,
			"001",
		},
	}

	for _, tt := range offsetTts {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := assert.New(t)
			ctx := context.Background()
			id := uuid.NewString()

			v1 := newTestAPIServerV1(t, id)
			defer tests.TearDown(t, v1.db, id)

			for i := 0; i < baseRecords; i++ {
				// Preserve semantic ordering from numbers.
				tests.StubUser(t, ctx, v1.db, fmt.Sprintf("%03d", i+1), model.UserRole_Student)
			}

			u := url.URL{Path: usersUrl}
			values := u.Query()
			values.Set("offset", tt.offset)
			u.RawQuery = values.Encode()

			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			resp, ok := v1.usersGet(req).(usersGetResponse)
			a.True(ok)

			if tt.wantUsers {
				a.Equal(tt.wantFirstID, resp.Users[0].ID)
			} else {
				a.Empty(resp.Users)
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
					Role: model.UserRole_Student,
				},
			},
			false,
			usersPostResponse{
				newSuccessResponse(),
				usersPostUserResponseFields{
					ID:   "NEW_USER",
					Role: model.UserRole_Student,
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
					Role: model.UserRole_Student,
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
					Role: model.UserRole_Student,
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
