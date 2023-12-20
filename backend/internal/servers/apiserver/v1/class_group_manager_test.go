package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
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

func TestAPIServerV1_classGroupManager(t *testing.T) {
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

			req := httpRequestWithAuthContext(
				httptest.NewRequest(tt.withMethod, fmt.Sprintf("%s%d", classGroupManagerUrl, 1), nil),
				tests.StubAuthContext(),
			)
			rr := httptest.NewRecorder()
			v1.classGroupManager(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_classGroupManagerGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withExistingClassGroupManager bool
		wantResponse                  classGroupManagerGetResponse
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with class group manager in database",
			true,
			classGroupManagerGetResponse{
				newSuccessResponse(),
				model.ClassGroupManager{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with class group manager not in database",
			false,
			classGroupManagerGetResponse{},
			http.StatusNotFound,
			"the requested class group manager does not exist",
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

			if tt.withExistingClassGroupManager {
				createdManager := tests.StubClassGroupManager(
					t, ctx, v1.db,
					tt.wantResponse.ClassGroupManager.ManagingRole,
					model.ClassType_Lec,
				)

				tt.wantResponse.ClassGroupManager.ID = createdManager.ID
				tt.wantResponse.ClassGroupManager.UserID = createdManager.UserID
				tt.wantResponse.ClassGroupManager.ClassGroupID = createdManager.ClassGroupID
				tt.wantResponse.ClassGroupManager.CreatedAt = createdManager.CreatedAt
				tt.wantResponse.ClassGroupManager.UpdatedAt = createdManager.CreatedAt
			}

			req := httpRequestWithAuthContext(
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%d", classGroupManagerUrl, tt.wantResponse.ClassGroupManager.ID), nil),
				tests.StubAuthContext(),
			)
			resp := v1.classGroupManagerGet(req, tt.wantResponse.ClassGroupManager.ID)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupManagerGetResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}

func TestAPIServerV1_classGroupManagerPatch(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withRequest                   classGroupManagerPatchRequest
		withExistingClassGroupManager bool
		wantResponse                  classGroupManagerPatchResponse
		wantNoChange                  bool
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with field changes",
			classGroupManagerPatchRequest{
				database.UpdateClassGroupManagerParams{
					ManagingRole: to.Ptr(model.ManagingRole_CourseCoordinator),
				},
			},
			true,
			classGroupManagerPatchResponse{
				newSuccessResponse(),
				classGroupManagerPatchClassGroupManagerResponseFields{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			false,
			http.StatusOK,
			"",
		},
		{
			"request with no field changes",
			classGroupManagerPatchRequest{
				database.UpdateClassGroupManagerParams{},
			},
			true,
			classGroupManagerPatchResponse{
				newSuccessResponse(),
				classGroupManagerPatchClassGroupManagerResponseFields{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			true,
			http.StatusOK,
			"",
		},
		{
			"request updating non-existent class group manager",
			classGroupManagerPatchRequest{
				database.UpdateClassGroupManagerParams{},
			},
			false,
			classGroupManagerPatchResponse{
				ClassGroupManager: classGroupManagerPatchClassGroupManagerResponseFields{
					ID: rand.Int63(),
				},
			},
			false,
			http.StatusNotFound,
			"class group manager to update does not exist",
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

			if tt.withExistingClassGroupManager {
				createdManager := tests.StubClassGroupManager(
					t, ctx, v1.db,
					tt.wantResponse.ClassGroupManager.ManagingRole,
					model.ClassType_Lec,
				)
				tt.wantResponse.ClassGroupManager.ID = createdManager.ID
				tt.wantResponse.ClassGroupManager.UserID = createdManager.UserID
				tt.wantResponse.ClassGroupManager.ClassGroupID = createdManager.ClassGroupID
				tt.wantResponse.ClassGroupManager.UpdatedAt = createdManager.CreatedAt
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			managerId := tt.wantResponse.ClassGroupManager.ID
			req := httpRequestWithAuthContext(
				httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", classGroupManagerUrl, managerId), bytes.NewReader(reqBodyBytes)),
				tests.StubAuthContext(),
			)
			resp := v1.classGroupManagerPatch(req, managerId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupManagerPatchResponse)
				a.True(ok)

				if !tt.wantNoChange {
					tt.wantResponse.ClassGroupManager.UpdatedAt = actualResp.ClassGroupManager.UpdatedAt
				}

				a.Equal(tt.wantResponse, actualResp)

				// Check that successive updates do not change anything.
				req = httpRequestWithAuthContext(
					httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", classGroupManagerUrl, managerId), bytes.NewReader(reqBodyBytes)),
					tests.StubAuthContext(),
				)
				successiveResp := v1.classGroupManagerPatch(req, managerId).(classGroupManagerPatchResponse)
				a.Equal(actualResp, successiveResp)
			}
		})
	}
}

func TestAPIServerV1_classGroupManagerDelete(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withExistingClassGroupManager bool
		wantResponse                  classGroupManagerDeleteResponse
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with existing class group manager",
			true,
			classGroupManagerDeleteResponse{newSuccessResponse()},
			http.StatusOK,
			"",
		},
		{
			"request with non-existent class group manager",
			false,
			classGroupManagerDeleteResponse{},
			http.StatusNotFound,
			"class group manager to delete does not exist",
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

			managerId := rand.Int63()
			if tt.withExistingClassGroupManager {
				createdManager := tests.StubClassGroupManager(
					t, ctx, v1.db,
					model.ManagingRole_CourseCoordinator,
					model.ClassType_Lec,
				)
				managerId = createdManager.ID
			}

			req := httpRequestWithAuthContext(
				httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s%d", classGroupManagerUrl, managerId), nil),
				tests.StubAuthContext(),
			)
			resp := v1.classGroupManagerDelete(req, managerId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupManagerDeleteResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
