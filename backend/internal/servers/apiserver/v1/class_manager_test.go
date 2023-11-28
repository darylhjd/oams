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
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/darylhjd/oams/backend/pkg/to"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_classManager(t *testing.T) {
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
			http.StatusNotImplemented,
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

			req := httptest.NewRequest(tt.withMethod, fmt.Sprintf("%s%d", classManagerUrl, 1), nil)
			rr := httptest.NewRecorder()
			v1.classManager(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_classManagerGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                     string
		withExistingClassManager bool
		wantResponse             classManagerGetResponse
		wantStatusCode           int
		wantErr                  string
	}{
		{
			"request with class manager in database",
			true,
			classManagerGetResponse{
				newSuccessResponse(),
				model.ClassManager{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with class manager not in database",
			false,
			classManagerGetResponse{},
			http.StatusNotFound,
			"the requested class manager does not exist",
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

			if tt.withExistingClassManager {
				createdManager := tests.StubClassManager(
					t, ctx, v1.db,
					tt.wantResponse.ClassManager.ManagingRole,
				)

				tt.wantResponse.ClassManager.ID = createdManager.ID
				tt.wantResponse.ClassManager.UserID = createdManager.UserID
				tt.wantResponse.ClassManager.ClassID = createdManager.ClassID
				tt.wantResponse.ClassManager.CreatedAt = createdManager.CreatedAt
				tt.wantResponse.ClassManager.UpdatedAt = createdManager.CreatedAt
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%d", classManagerUrl, tt.wantResponse.ClassManager.ID), nil)
			resp := v1.classManagerGet(req, tt.wantResponse.ClassManager.ID)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classManagerGetResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}

func TestAPIServerV1_classManagerPatch(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                     string
		withRequest              classManagerPatchRequest
		withExistingClassManager bool
		wantResponse             classManagerPatchResponse
		wantNoChange             bool
		wantStatusCode           int
		wantErr                  string
	}{
		{
			"request with field changes",
			classManagerPatchRequest{
				database.UpdateClassManagerParams{
					ManagingRole: to.Ptr(model.ManagingRole_CourseCoordinator),
				},
			},
			true,
			classManagerPatchResponse{
				newSuccessResponse(),
				classManagerPatchClassManagerResponseFields{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			false,
			http.StatusOK,
			"",
		},
		{
			"request with no field changes",
			classManagerPatchRequest{
				database.UpdateClassManagerParams{},
			},
			true,
			classManagerPatchResponse{
				newSuccessResponse(),
				classManagerPatchClassManagerResponseFields{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			true,
			http.StatusOK,
			"",
		},
		{
			"request updating non-existent class manager",
			classManagerPatchRequest{
				database.UpdateClassManagerParams{},
			},
			false,
			classManagerPatchResponse{
				ClassManager: classManagerPatchClassManagerResponseFields{
					ID: 6666,
				},
			},
			false,
			http.StatusNotFound,
			"class manager to update does not exist",
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

			if tt.withExistingClassManager {
				createdManager := tests.StubClassManager(t, ctx, v1.db, tt.wantResponse.ClassManager.ManagingRole)
				tt.wantResponse.ClassManager.ID = createdManager.ID
				tt.wantResponse.ClassManager.UserID = createdManager.UserID
				tt.wantResponse.ClassManager.ClassID = createdManager.ClassID
				tt.wantResponse.ClassManager.UpdatedAt = createdManager.CreatedAt
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			managerId := tt.wantResponse.ClassManager.ID
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", classManagerUrl, managerId), bytes.NewReader(reqBodyBytes))
			resp := v1.classManagerPatch(req, managerId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classManagerPatchResponse)
				a.True(ok)

				if !tt.wantNoChange {
					tt.wantResponse.ClassManager.UpdatedAt = actualResp.ClassManager.UpdatedAt
				}

				a.Equal(tt.wantResponse, actualResp)

				// Check that successive updates do not change anything.
				req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", classManagerUrl, managerId), bytes.NewReader(reqBodyBytes))
				successiveResp := v1.classManagerPatch(req, managerId).(classManagerPatchResponse)
				a.Equal(actualResp, successiveResp)
			}
		})
	}
}
