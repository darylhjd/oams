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
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/darylhjd/oams/backend/pkg/to"
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
			"with PATCH method",
			http.MethodPatch,
			http.StatusBadRequest,
		},
		{
			"with DELETE method",
			http.MethodDelete,
			http.StatusNotFound,
		},
		{
			"with PUT method",
			http.MethodPut,
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

func TestAPIServerV1_userMe(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name             string
		withStubAuthUser bool
		wantResponse     userMeResponse
		wantStatusCode   int
		wantErr          string
	}{
		{
			"request with valid auth context and auth context user in database",
			true,
			userMeResponse{
				newSuccessResponse(),
				tests.StubAuthContext().User,
				false,
			},
			http.StatusOK,
			"",
		},
		{
			"request with valid auth context but non-existent user in database",
			false,
			userMeResponse{},
			http.StatusInternalServerError,
			"could get session user from database",
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

			if tt.withStubAuthUser {
				tests.StubAuthContextUser(t, ctx, v1.db)
			}

			req := httpRequestWithAuthContext(
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", userUrl, sessionUserId), nil),
				tests.StubAuthContext(),
			)
			resp := v1.userMe(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(userMeResponse)
				a.True(ok)

				tt.wantResponse.SessionUser.CreatedAt = actualResp.SessionUser.CreatedAt
				tt.wantResponse.SessionUser.UpdatedAt = actualResp.SessionUser.UpdatedAt
				a.Equal(tt.wantResponse, actualResp)
			}
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
				model.User{
					ID:   uuid.NewString(),
					Role: model.UserRole_User,
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
				createdUser := tests.StubUser(t, ctx, v1.db, database.CreateUserParams{
					ID:   tt.wantResponse.User.ID,
					Role: tt.wantResponse.User.Role,
				})
				tt.wantResponse.User.CreatedAt = createdUser.CreatedAt
				tt.wantResponse.User.UpdatedAt = createdUser.CreatedAt
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

func TestAPIServerV1_userPatch(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name             string
		withRequest      userPatchRequest
		withExistingUser bool
		wantResponse     userPatchResponse
		wantNoChange     bool
		wantStatusCode   int
		wantErr          string
	}{
		{
			"request with field changes",
			userPatchRequest{
				database.UpdateUserParams{
					Name:  to.Ptr("NEW NAME"),
					Email: to.Ptr("NEW EMAIL"),
				},
			},
			true,
			userPatchResponse{
				newSuccessResponse(),
				userPatchUserResponseFields{
					ID:    uuid.NewString(),
					Name:  "NEW NAME",
					Email: "NEW EMAIL",
					Role:  model.UserRole_User,
				},
			},
			false,
			http.StatusOK,
			"",
		},
		{
			"request with no field changes",
			userPatchRequest{
				database.UpdateUserParams{},
			},
			true,
			userPatchResponse{
				newSuccessResponse(),
				userPatchUserResponseFields{
					ID:   uuid.NewString(),
					Role: model.UserRole_User,
				},
			},
			true,
			http.StatusOK,
			"",
		},
		{
			"request updating non-existent user",
			userPatchRequest{
				database.UpdateUserParams{},
			},
			false,
			userPatchResponse{},
			false,
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
				createdUser := tests.StubUser(t, ctx, v1.db, database.CreateUserParams{
					ID:   tt.wantResponse.User.ID,
					Role: tt.wantResponse.User.Role,
				})
				tt.wantResponse.User.UpdatedAt = createdUser.CreatedAt
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%s", userUrl, userId), bytes.NewReader(reqBodyBytes))
			resp := v1.userPatch(req, userId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(userPatchResponse)
				a.True(ok)

				if !tt.wantNoChange {
					tt.wantResponse.User.UpdatedAt = actualResp.User.UpdatedAt
				}

				a.Equal(tt.wantResponse, actualResp)

				// Check that successive updates do not change anything.
				req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%s", userUrl, userId), bytes.NewReader(reqBodyBytes))
				successiveResp := v1.userPatch(req, userId).(userPatchResponse)
				a.Equal(actualResp, successiveResp)
			}
		})
	}
}

func TestAPIServerV1_userDelete(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                     string
		withExistingUser         bool
		withForeignKeyDependency bool
		wantResponse             userDeleteResponse
		wantStatusCode           int
		wantErr                  string
	}{
		{
			"request with existing user",
			true,
			false,
			userDeleteResponse{newSuccessResponse()},
			http.StatusOK,
			"",
		},
		{
			"request with non-existent user",
			false,
			false,
			userDeleteResponse{},
			http.StatusNotFound,
			"user to delete does not exist",
		},
		{
			"request with user foreign key constraint",
			true,
			true,
			userDeleteResponse{},
			http.StatusConflict,
			"user to delete is still referenced",
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

			var userId string
			switch {
			case tt.withForeignKeyDependency:
				createdEnrollment := tests.StubSessionEnrollment(t, ctx, v1.db, true)
				userId = createdEnrollment.UserID
			case tt.withExistingUser:
				userId = tests.StubUser(t, ctx, v1.db, database.CreateUserParams{
					ID:   uuid.NewString(),
					Role: model.UserRole_User,
				}).ID
			default:
				userId = uuid.NewString()
			}

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s", userUrl, userId), nil)
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
