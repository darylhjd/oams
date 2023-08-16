package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
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
			http.StatusBadRequest,
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
						Name:      "NEW21",
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

func TestAPIServerV1_classGroupsPost(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                   string
		withRequest            classGroupsPostRequest
		withExistingClassGroup bool
		withExistingClass      bool
		wantResponse           classGroupsPostResponse
		wantStatusCode         int
		wantErr                string
	}{
		{
			"request with no existing class group",
			classGroupsPostRequest{
				database.CreateClassGroupParams{
					Name:      "NEW21",
					ClassType: model.ClassType_Lab,
				},
			},
			false,
			true,
			classGroupsPostResponse{
				newSuccessResponse(),
				classGroupsPostClassGroupResponseFields{
					Name:      "NEW21",
					ClassType: model.ClassType_Lab,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with existing class group",
			classGroupsPostRequest{
				database.CreateClassGroupParams{
					Name:      "EXISTING22",
					ClassType: model.ClassType_Lec,
				},
			},
			true,
			true,
			classGroupsPostResponse{},
			http.StatusConflict,
			"class group with same class_id, name, and class_type already exists",
		},
		{
			"request with non-existent class dependency",
			classGroupsPostRequest{
				database.CreateClassGroupParams{
					ClassID:   rand.Int63(),
					Name:      "FAIL_INSERT22",
					ClassType: model.ClassType_Tut,
				},
			},
			false,
			false,
			classGroupsPostResponse{},
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

			switch {
			case tt.withExistingClassGroup:
				createdGroup := tests.StubClassGroup(
					t, ctx, v1.db,
					tt.withRequest.ClassGroup.Name,
					tt.withRequest.ClassGroup.ClassType,
				)
				tt.withRequest.ClassGroup.ClassID = createdGroup.ClassID
			case tt.withExistingClass:
				createdClass := tests.StubClass(
					t, ctx, v1.db,
					uuid.NewString(),
					rand.Int31(),
					uuid.NewString(),
				)
				tt.withRequest.ClassGroup.ClassID = createdClass.ID
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPost, classGroupsUrl, bytes.NewReader(reqBodyBytes))
			resp := v1.classGroupsPost(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupsPostResponse)
				a.True(ok)

				tt.wantResponse.ClassGroup.ID = actualResp.ClassGroup.ID
				tt.wantResponse.ClassGroup.ClassID = actualResp.ClassGroup.ClassID
				tt.wantResponse.ClassGroup.CreatedAt = actualResp.ClassGroup.CreatedAt
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
