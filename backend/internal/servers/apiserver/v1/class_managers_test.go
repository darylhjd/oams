package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
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

func TestAPIServerV1_classManagersGetQueryParams(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name           string
		query          url.Values
		wantStatusCode int
		wantErr        string
	}{
		{
			"sort with correct column",
			url.Values{
				"sort": []string{"user_id"},
			},
			http.StatusOK,
			"",
		},
		{
			"sort with wrong column",
			url.Values{
				"sort": []string{"wrong"},
			},
			http.StatusBadRequest,
			"unknown sort column `wrong`",
		},
		{
			"sort with no value",
			url.Values{
				"sort": []string{},
			},
			http.StatusOK,
			"",
		},
		{
			"limit present",
			url.Values{
				"limit": []string{"1"},
			},
			http.StatusOK,
			"",
		},
		{
			"limit with no value",
			url.Values{
				"limit": []string{},
			},
			http.StatusOK,
			"",
		},
		{
			"offset present",
			url.Values{
				"offset": []string{"1"},
			},
			http.StatusOK,
			"",
		},
		{
			"offset with no value",
			url.Values{
				"offset": []string{},
			},
			http.StatusOK,
			"",
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

			u := url.URL{Path: classManagersUrl}
			u.RawQuery = tt.query.Encode()

			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			resp := v1.classManagersGet(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				_, ok := resp.(classManagersGetResponse)
				a.True(ok)
			}
		})
	}
}

func TestAPIServerV1_classManagersPost(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                     string
		withRequest              classManagersPostRequest
		withExistingClassManager bool
		withExistingUser         bool
		withExistingClass        bool
		wantResponse             classManagersPostResponse
		wantStatusCode           int
		wantErr                  string
	}{
		{
			"request with no existing class manager",
			classManagersPostRequest{
				database.CreateClassManagerParams{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			false,
			true,
			true,
			classManagersPostResponse{
				newSuccessResponse(),
				classManagersPostClassManagerResponseFields{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with existing class manager",
			classManagersPostRequest{
				database.CreateClassManagerParams{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			true,
			true,
			true,
			classManagersPostResponse{},
			http.StatusConflict,
			"class manager with same user_id and class_id already exists",
		},
		{
			"request with non-existent user dependency",
			classManagersPostRequest{
				database.CreateClassManagerParams{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			false,
			false,
			true,
			classManagersPostResponse{},
			http.StatusBadRequest,
			"user_id and/or class_id does not exist",
		},
		{
			"request with non-existent class dependency",
			classManagersPostRequest{
				database.CreateClassManagerParams{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			false,
			true,
			false,
			classManagersPostResponse{},
			http.StatusBadRequest,
			"user_id and/or class_id does not exist",
		},
		{
			"request with non-existent user and class dependency",
			classManagersPostRequest{
				database.CreateClassManagerParams{
					ManagingRole: model.ManagingRole_CourseCoordinator,
				},
			},
			false,
			false,
			false,
			classManagersPostResponse{},
			http.StatusBadRequest,
			"user_id and/or class_id does not exist",
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
			case tt.withExistingClassManager:
				createdManager := tests.StubClassManager(
					t, ctx, v1.db,
					tt.withRequest.ClassManager.ManagingRole,
				)
				tt.withRequest.ClassManager.UserID = createdManager.UserID
				tt.withRequest.ClassManager.ClassID = createdManager.ClassID
			default:
				if tt.withExistingUser {
					createdUser := tests.StubUser(t, ctx, v1.db, uuid.NewString(), model.UserRole_User)
					tt.withRequest.ClassManager.UserID = createdUser.ID
				}

				if tt.withExistingClass {
					createdClass := tests.StubClass(
						t, ctx, v1.db,
						uuid.NewString(),
						rand.Int31(),
						uuid.NewString(),
					)
					tt.withRequest.ClassManager.ClassID = createdClass.ID
				}
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPost, classManagersUrl, bytes.NewReader(reqBodyBytes))
			resp := v1.classManagersPost(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classManagersPostResponse)
				a.True(ok)

				tt.wantResponse.ClassManager.ID = actualResp.ClassManager.ID
				tt.wantResponse.ClassManager.UserID = actualResp.ClassManager.UserID
				tt.wantResponse.ClassManager.ClassID = actualResp.ClassManager.ClassID
				tt.wantResponse.ClassManager.CreatedAt = actualResp.ClassManager.CreatedAt
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
