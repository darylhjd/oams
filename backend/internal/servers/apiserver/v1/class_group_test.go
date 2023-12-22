package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_classGroup(t *testing.T) {
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
			http.StatusNotImplemented,
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

			req := httptest.NewRequest(tt.withMethod, fmt.Sprintf("%s%d", classGroupUrl, 1), nil)
			rr := httptest.NewRecorder()
			v1.classGroup(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_classGroupGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                   string
		withExistingClassGroup bool
		wantResponse           classGroupGetResponse
		wantStatusCode         int
		wantErr                string
	}{
		{
			"request with class group in database",
			true,
			classGroupGetResponse{
				newSuccessResponse(),
				model.ClassGroup{
					Name:      uuid.NewString(),
					ClassType: model.ClassType_Lec,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with class group not in database",
			false,
			classGroupGetResponse{},
			http.StatusNotFound,
			"the requested class group does not exist",
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
				createdClassGroup := tests.StubClassGroup(
					t, ctx, v1.db,
					tt.wantResponse.ClassGroup.Name,
					tt.wantResponse.ClassGroup.ClassType,
				)

				tt.wantResponse.ClassGroup.ID = createdClassGroup.ID
				tt.wantResponse.ClassGroup.ClassID = createdClassGroup.ClassID
				tt.wantResponse.ClassGroup.CreatedAt = createdClassGroup.CreatedAt
				tt.wantResponse.ClassGroup.UpdatedAt = createdClassGroup.CreatedAt
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%d", classGroupUrl, tt.wantResponse.ClassGroup.ID), nil)
			resp := v1.classGroupGet(req, tt.wantResponse.ClassGroup.ID)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupGetResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
