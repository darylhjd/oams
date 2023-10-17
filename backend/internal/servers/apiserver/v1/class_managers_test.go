package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_classManagers(t *testing.T) {
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

			req := httptest.NewRequest(tt.withMethod, classManagersUrl, nil)
			rr := httptest.NewRecorder()
			v1.classManagers(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_classManagersGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                     string
		withExistingClassManager bool
		wantResponse             classManagersGetResponse
	}{
		{
			"request with class manager in database",
			true,
			classManagersGetResponse{
				newSuccessResponse(),
				[]model.ClassManager{
					{
						ManagingRole: model.ManagingRole_CourseCoordinator,
					},
				},
			},
		},
		{
			"request with no class manager in database",
			false,
			classManagersGetResponse{
				newSuccessResponse(),
				[]model.ClassManager{},
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

			if tt.withExistingClassManager {
				for idx, manager := range tt.wantResponse.ClassManagers {
					createdManager := tests.StubClassManager(t, ctx, v1.db, manager.ManagingRole)
					managerPtr := &tt.wantResponse.ClassManagers[idx]
					managerPtr.ID = createdManager.ID
					managerPtr.UserID = createdManager.UserID
					managerPtr.ClassID = createdManager.ClassID
					managerPtr.CreatedAt, managerPtr.UpdatedAt = createdManager.CreatedAt, createdManager.CreatedAt
				}
			}

			req := httptest.NewRequest(http.MethodGet, classManagersUrl, nil)
			actualResp, ok := v1.classManagersGet(req).(classManagersGetResponse)
			a.True(ok)
			a.Equal(tt.wantResponse, actualResp)
		})
	}
}
