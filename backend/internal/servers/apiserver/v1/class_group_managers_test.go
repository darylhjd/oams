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
			http.StatusBadRequest,
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

func TestAPIServerV1_classGroupManagersPost(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withRequest                   classGroupManagersPostRequest
		withExistingClassGroupManager bool
		withExistingUser              bool
		withExistingClass             bool
		wantResponse                  classGroupManagersPostResponse
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with no existing class group manager",
			classGroupManagersPostRequest{
				database.CreateClassGroupManagerParams{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			false,
			true,
			true,
			classGroupManagersPostResponse{
				newSuccessResponse(),
				classGroupManagersPostClassGroupManagerResponseFields{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with existing class group manager",
			classGroupManagersPostRequest{
				database.CreateClassGroupManagerParams{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			true,
			true,
			true,
			classGroupManagersPostResponse{},
			http.StatusConflict,
			"class group manager with same user_id and class_group_id already exists",
		},
		{
			"request with non-existent user dependency",
			classGroupManagersPostRequest{
				database.CreateClassGroupManagerParams{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			false,
			false,
			true,
			classGroupManagersPostResponse{},
			http.StatusBadRequest,
			"user_id and/or class_group_id does not exist",
		},
		{
			"request with non-existent class dependency",
			classGroupManagersPostRequest{
				database.CreateClassGroupManagerParams{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			false,
			true,
			false,
			classGroupManagersPostResponse{},
			http.StatusBadRequest,
			"user_id and/or class_group_id does not exist",
		},
		{
			"request with non-existent user and class dependency",
			classGroupManagersPostRequest{
				database.CreateClassGroupManagerParams{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			false,
			false,
			false,
			classGroupManagersPostResponse{},
			http.StatusBadRequest,
			"user_id and/or class_group_id does not exist",
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
			case tt.withExistingClassGroupManager:
				createdManager := tests.StubClassGroupManager(
					t, ctx, v1.db,
					tt.withRequest.ClassGroupManager.ManagingRole,
					model.ClassType_Lec,
				)
				tt.withRequest.ClassGroupManager.UserID = createdManager.UserID
				tt.withRequest.ClassGroupManager.ClassGroupID = createdManager.ClassGroupID
			default:
				if tt.withExistingUser {
					createdUser := tests.StubUser(t, ctx, v1.db, database.CreateUserParams{
						ID:    uuid.NewString(),
						Name:  uuid.NewString(),
						Email: uuid.NewString(),
						Role:  model.UserRole_SystemAdmin,
					})
					tt.withRequest.ClassGroupManager.UserID = createdUser.ID
				}

				if tt.withExistingClass {
					createdClass := tests.StubClass(t, ctx, v1.db, database.CreateClassParams{
						Code:      uuid.NewString(),
						Year:      rand.Int31(),
						Semester:  uuid.NewString(),
						Programme: uuid.NewString(),
					})
					tt.withRequest.ClassGroupManager.ClassGroupID = createdClass.ID
				}
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPost, classGroupManagersUrl, bytes.NewReader(reqBodyBytes))
			resp := v1.classGroupManagersPost(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupManagersPostResponse)
				a.True(ok)

				tt.wantResponse.ClassGroupManager.ID = actualResp.ClassGroupManager.ID
				tt.wantResponse.ClassGroupManager.UserID = actualResp.ClassGroupManager.UserID
				tt.wantResponse.ClassGroupManager.ClassGroupID = actualResp.ClassGroupManager.ClassGroupID
				tt.wantResponse.ClassGroupManager.CreatedAt = actualResp.ClassGroupManager.CreatedAt
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
