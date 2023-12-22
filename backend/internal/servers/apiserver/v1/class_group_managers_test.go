package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_classGroupManagers(t *testing.T) {
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
			http.StatusNotImplemented,
		},
		{
			"with PUT method",
			http.MethodPut,
			http.StatusNotImplemented,
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

			req := httpRequestWithAuthContext(
				httptest.NewRequest(tt.withMethod, classGroupManagersUrl, nil),
				tests.StubAuthContext(),
			)
			rr := httptest.NewRecorder()
			v1.classGroupManagers(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_classGroupManagersGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withExistingClassGroupManager bool
		wantResponse                  classGroupManagersGetResponse
	}{
		{
			"request with class group manager in database",
			true,
			classGroupManagersGetResponse{
				newSuccessResponse(),
				[]model.ClassGroupManager{
					{
						ManagingRole: model.ManagingRole_CourseCoordinator,
					},
				},
			},
		},
		{
			"request with no class group manager in database",
			false,
			classGroupManagersGetResponse{
				newSuccessResponse(),
				[]model.ClassGroupManager{},
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

			if tt.withExistingClassGroupManager {
				for idx, manager := range tt.wantResponse.ClassGroupManagers {
					createdManager := tests.StubClassGroupManager(
						t, ctx, v1.db,
						manager.ManagingRole,
						model.ClassType_Lec,
					)
					managerPtr := &tt.wantResponse.ClassGroupManagers[idx]
					managerPtr.ID = createdManager.ID
					managerPtr.UserID = createdManager.UserID
					managerPtr.ClassGroupID = createdManager.ClassGroupID
					managerPtr.CreatedAt, managerPtr.UpdatedAt = createdManager.CreatedAt, createdManager.CreatedAt
				}
			}

			req := httpRequestWithAuthContext(
				httptest.NewRequest(http.MethodGet, classGroupManagersUrl, nil),
				tests.StubAuthContext(),
			)
			actualResp, ok := v1.classGroupManagersGet(req).(classGroupManagersGetResponse)
			a.True(ok)
			a.Equal(tt.wantResponse, actualResp)
		})
	}
}
