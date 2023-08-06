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
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func TestAPIServerV1_classGroupPatch(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                    string
		withRequest             classGroupPatchRequest
		withExistingClassGroup  bool
		withUpdateConflict      bool
		withExistingUpdateClass bool
		wantResponse            classGroupPatchResponse
		wantNoChange            bool
		wantStatusCode          int
		wantErr                 string
	}{
		{
			"request with field changes",
			classGroupPatchRequest{
				classGroupPatchClassGroupRequestFields{
					Name:      ptr("NEW21"),
					ClassType: ptr(database.ClassTypeLAB),
				},
			},
			true,
			false,
			true,
			classGroupPatchResponse{
				newSuccessResponse(),
				database.UpdateClassGroupRow{
					Name:      "NEW21",
					ClassType: database.ClassTypeLAB,
				},
			},
			false,
			http.StatusOK,
			"",
		},
		{
			"request with no field changes",
			classGroupPatchRequest{
				classGroupPatchClassGroupRequestFields{},
			},
			true,
			false,
			true,
			classGroupPatchResponse{
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
			classGroupPatchRequest{
				classGroupPatchClassGroupRequestFields{},
			},
			false,
			false,
			false,
			classGroupPatchResponse{},
			false,
			http.StatusNotFound,
			"class group to update does not exist",
		},
		{
			"request with update conflict",
			classGroupPatchRequest{
				classGroupPatchClassGroupRequestFields{
					Name:      ptr("EXISTING32"),
					ClassType: ptr(database.ClassTypeLAB),
				},
			},
			true,
			true,
			true,
			classGroupPatchResponse{},
			false,
			http.StatusConflict,
			"class group with same class_id and name already exists",
		},
		{
			"request with non existent class dependency",
			classGroupPatchRequest{},
			true,
			false,
			false,
			classGroupPatchResponse{},
			false,
			http.StatusBadRequest,
			"class_id does not exist",
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

			var groupId int64
			switch {
			case tt.withUpdateConflict:
				// Create group to update.
				updateClassGroup := tests.StubClassGroup(t, ctx, v1.db.Q, uuid.NewString(), database.ClassTypeTUT)
				groupId = updateClassGroup.ID

				// Also create group to conflict with.
				_ = tests.StubClassGroupWithClassID(
					t, ctx, v1.db.Q,
					updateClassGroup.ClassID,
					*tt.withRequest.ClassGroup.Name,
					*tt.withRequest.ClassGroup.ClassType,
				)
			case tt.withExistingClassGroup && !tt.withExistingUpdateClass:
				createdClassGroup := tests.StubClassGroup(
					t, ctx, v1.db.Q,
					uuid.NewString(),
					database.ClassTypeTUT,
				)

				groupId = createdClassGroup.ID
				tt.withRequest.ClassGroup.ClassID = ptr(createdClassGroup.ClassID + 1)
			case tt.withExistingClassGroup:
				createdClassGroup := tests.StubClassGroup(
					t, ctx, v1.db.Q,
					tt.wantResponse.ClassGroup.Name,
					tt.wantResponse.ClassGroup.ClassType,
				)

				groupId = createdClassGroup.ID
				tt.wantResponse.ClassGroup.ID = createdClassGroup.ID
				tt.wantResponse.ClassGroup.ClassID = createdClassGroup.ClassID
				tt.wantResponse.ClassGroup.UpdatedAt = createdClassGroup.CreatedAt
			default:
				groupId = rand.Int63()
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", classGroupUrl, groupId), bytes.NewReader(reqBodyBytes))
			resp := v1.classGroupPatch(req, groupId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupPatchResponse)
				a.True(ok)

				if !tt.wantNoChange {
					tt.wantResponse.ClassGroup.UpdatedAt = actualResp.ClassGroup.UpdatedAt
				}

				a.Equal(tt.wantResponse, actualResp)

				// Check that successive updates do not change anything.
				req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", classGroupUrl, groupId), bytes.NewReader(reqBodyBytes))
				successiveResp := v1.classGroupPatch(req, groupId).(classGroupPatchResponse)
				a.Equal(actualResp, successiveResp)
			}
		})
	}
}

func TestAPIServerV1_classGroupDelete(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                     string
		withExistingClassGroup   bool
		withForeignKeyDependency bool
		wantResponse             classGroupDeleteResponse
		wantStatusCode           int
		wantErr                  string
	}{
		{
			"request with existing class group",
			true,
			false,
			classGroupDeleteResponse{newSuccessResponse()},
			http.StatusOK,
			"",
		},
		{
			"request with non-existent class group",
			false,
			false,
			classGroupDeleteResponse{},
			http.StatusNotFound,
			"class group to delete does not exist",
		},
		{
			"request with class group foreign key dependency",
			true,
			true,
			classGroupDeleteResponse{},
			http.StatusConflict,
			"class group to delete is still referenced",
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

			var groupId int64
			switch {
			case tt.withForeignKeyDependency:
				createdClassGroupSession := tests.StubClassGroupSession(
					t, ctx, v1.db.Q,
					pgtype.Timestamptz{Time: time.UnixMicro(1), Valid: true},
					pgtype.Timestamptz{Time: time.UnixMicro(2), Valid: true},
					uuid.NewString(),
				)
				groupId = createdClassGroupSession.ClassGroupID
			case tt.withExistingClassGroup:
				createdClassGroup := tests.StubClassGroup(
					t, ctx, v1.db.Q,
					uuid.NewString(),
					database.ClassTypeTUT,
				)
				groupId = createdClassGroup.ID
			default:
				groupId = rand.Int63()
			}

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s%d", classGroupUrl, groupId), nil)
			resp := v1.classGroupDelete(req, groupId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupDeleteResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
