package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_user(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name           string
		withMethod     string
		wantStatusCode int
	}{
		{
			"with GET method",
			http.MethodGet,
			http.StatusNotFound,
		},
		{
			"with PUT method",
			http.MethodPut,
			http.StatusBadRequest,
		},
		{
			"with DELETE method",
			http.MethodDelete,
			http.StatusNotFound,
		},
		{
			"with PATCH method",
			http.MethodPatch,
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

			req := httptest.NewRequest(tt.withMethod, userUrl, nil)
			rr := httptest.NewRecorder()
			v1.user(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_userGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name             string
		withExistingUser bool
		wantResponse     userGetResponse
		wantStatusCode   int
		wantErr          string
	}{
		{
			"request with user in database",
			true,
			userGetResponse{
				newSuccessResponse(),
				userGetUserResponseFields{
					ID:   "EXISTING_USER",
					Role: database.UserRoleSTUDENT,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with user not in database",
			false,
			userGetResponse{},
			http.StatusNotFound,
			"the requested user does not exist",
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
				tests.StubUser(t, ctx, v1.db.Q, tt.wantResponse.User.ID)
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", userUrl, tt.wantResponse.User.ID), nil)
			resp := v1.userGet(req, tt.wantResponse.User.ID)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(userGetResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}

func TestAPIServerV1_userPut(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name             string
		withRequest      userPutRequest
		withExistingUser bool
		wantResponse     userPutResponse
		wantStatusCode   int
		wantErr          string
	}{
		{
			"request with all fields set",
			userPutRequest{
				userPutUserRequestFields{
					ptr("NEW NAME"),
					ptr("NEW EMAIL"),
					ptr(database.UserRoleSTUDENT),
				},
			},
			true,
			userPutResponse{
				newSuccessResponse(),
				userPutUserResponseFields{
					ID:    "NEW_ID",
					Name:  "NEW NAME",
					Email: ptr("NEW EMAIL"),
					Role:  database.UserRoleSTUDENT,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with optional fields not set",
			userPutRequest{
				userPutUserRequestFields{},
			},
			true,
			userPutResponse{
				newSuccessResponse(),
				userPutUserResponseFields{
					ID:   "NEW_ID",
					Role: database.UserRoleSTUDENT,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request updating non-existent user",
			userPutRequest{
				userPutUserRequestFields{},
			},
			false,
			userPutResponse{},
			http.StatusNotFound,
			"user to update does not exist",
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

			userId := tt.wantResponse.User.ID
			if tt.withExistingUser {
				tests.StubUser(t, ctx, v1.db.Q, userId)
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s%s", userUrl, userId), bytes.NewReader(reqBodyBytes))
			resp := v1.userPut(req, userId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(userPutResponse)
				a.True(ok)

				tt.wantResponse.User.UpdatedAt = actualResp.User.UpdatedAt
				a.Equal(tt.wantResponse, actualResp)

				// Check that successive updates do not change the updated_at field.
				req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s%s", userUrl, userId), bytes.NewReader(reqBodyBytes))
				successiveResp := v1.userPut(req, userId).(userPutResponse)
				a.Equal(actualResp.User.UpdatedAt, successiveResp.User.UpdatedAt)
			}
		})
	}
}

func TestAPIServerV1_userDelete(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name             string
		withExistingUser bool
		wantResponse     userDeleteResponse
		wantStatusCode   int
		wantErr          string
	}{
		{
			"request with existing user",
			true,
			userDeleteResponse{newSuccessResponse()},
			http.StatusOK,
			"",
		},
		{
			"request with non-existent user",
			false,
			userDeleteResponse{},
			http.StatusNotFound,
			"user to delete does not exist",
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

			userId := uuid.NewString()
			if tt.withExistingUser {
				tests.StubUser(t, ctx, v1.db.Q, userId)
			}

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s%s", userUrl, userId), nil)
			resp := v1.userDelete(req, userId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(userDeleteResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
