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
			http.StatusNotImplemented,
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
					"NEW_ID",
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
				userPutUserRequestFields{
					ID: "NEW_ID",
				},
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
				userPutUserRequestFields{
					ID: "NON_EXISTENT_USER",
				},
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

			if tt.withExistingUser {
				tests.StubUser(t, ctx, v1.db.Q, tt.withRequest.User.ID)
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPut, userUrl, bytes.NewReader(reqBodyBytes))
			resp := v1.userPut(req)
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
			}
		})
	}
}
