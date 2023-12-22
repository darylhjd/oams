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

func TestAPIServerV1_classGroups(t *testing.T) {
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

			req := httptest.NewRequest(tt.withMethod, classGroupsUrl, nil)
			rr := httptest.NewRecorder()
			v1.classGroups(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_classGroupsGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                   string
		withExistingClassGroup bool
		wantResponse           classGroupsGetResponse
	}{
		{
			"request with class group in database",
			true,
			classGroupsGetResponse{
				newSuccessResponse(),
				[]model.ClassGroup{
					{
						Name:      uuid.NewString(),
						ClassType: model.ClassType_Lec,
					},
				},
			},
		},
		{
			"request with no class group in database",
			false,
			classGroupsGetResponse{
				newSuccessResponse(),
				[]model.ClassGroup{},
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

			if tt.withExistingClassGroup {
				for idx, group := range tt.wantResponse.ClassGroups {
					createdGroup := tests.StubClassGroup(t, ctx, v1.db, group.Name, group.ClassType)
					groupPtr := &tt.wantResponse.ClassGroups[idx]
					groupPtr.ID = createdGroup.ID
					groupPtr.ClassID = createdGroup.ClassID
					groupPtr.CreatedAt, groupPtr.UpdatedAt = createdGroup.CreatedAt, createdGroup.CreatedAt
				}
			}

			req := httptest.NewRequest(http.MethodGet, classGroupsUrl, nil)
			actualResp, ok := v1.classGroupsGet(req).(classGroupsGetResponse)
			a.True(ok)
			a.Equal(tt.wantResponse, actualResp)
		})
	}
}
