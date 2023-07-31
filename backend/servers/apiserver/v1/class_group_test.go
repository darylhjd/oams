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
			"with PUT method",
			http.MethodPut,
			http.StatusNotImplemented,
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
				database.ClassGroup{
					Name:      "EXISTING21",
					ClassType: database.ClassTypeLEC,
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
					t, ctx, v1.db.Q,
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

func TestAPIServerV1_classGroupPut(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                   string
		withRequest            classGroupPutRequest
		withExistingClassGroup bool
		wantResponse           classGroupPutResponse
		wantNoChange           bool
		wantStatusCode         int
		wantErr                string
	}{
		{
			"request with all fields set",
			classGroupPutRequest{
				classGroupPutClassGroupRequestFields{
					ptr(int64(1)),
					ptr("NEW21"),
					ptr(database.ClassTypeLAB),
				},
			},
			true,
			classGroupPutResponse{
				newSuccessResponse(),
				database.UpdateClassGroupRow{
					ClassID:   1,
					Name:      "NEW21",
					ClassType: database.ClassTypeLAB,
				},
			},
			false,
			http.StatusOK,
			"",
		},
		{
			"request with optional fields not set",
			classGroupPutRequest{
				classGroupPutClassGroupRequestFields{},
			},
			true,
			classGroupPutResponse{
				newSuccessResponse(),
				database.UpdateClassGroupRow{
					Name:      "EXISTING21",
					ClassType: database.ClassTypeLEC,
				},
			},
			true,
			http.StatusOK,
			"",
		},
		{
			"request updating non-existent class group",
			classGroupPutRequest{
				classGroupPutClassGroupRequestFields{},
			},
			false,
			classGroupPutResponse{
				ClassGroup: database.UpdateClassGroupRow{
					ID: 6666,
				},
			},
			false,
			http.StatusNotFound,
			"class group to update does not exist",
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
					t, ctx, v1.db.Q,
					tt.wantResponse.ClassGroup.Name,
					tt.wantResponse.ClassGroup.ClassType,
				)
				tt.wantResponse.ClassGroup.ID = createdClassGroup.ID
				tt.wantResponse.ClassGroup.ClassID = createdClassGroup.ClassID
				tt.wantResponse.ClassGroup.UpdatedAt = createdClassGroup.CreatedAt
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			groupId := tt.wantResponse.ClassGroup.ID
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s%d", classGroupUrl, groupId), bytes.NewReader(reqBodyBytes))
			resp := v1.classGroupPut(req, groupId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupPutResponse)
				a.True(ok)

				if !tt.wantNoChange {
					tt.wantResponse.ClassGroup.UpdatedAt = actualResp.ClassGroup.UpdatedAt
				}

				a.Equal(tt.wantResponse, actualResp)

				// Check that successive updates do not change anything.
				req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s%d", classGroupUrl, groupId), bytes.NewReader(reqBodyBytes))
				successiveResp := v1.classGroupPut(req, groupId).(classGroupPutResponse)
				a.Equal(actualResp, successiveResp)
			}
		})
	}
}
