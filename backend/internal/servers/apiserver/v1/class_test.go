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

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/darylhjd/oams/backend/pkg/to"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_class(t *testing.T) {
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

			req := httptest.NewRequest(tt.withMethod, fmt.Sprintf("%s%d", classUrl, 1), nil)
			rr := httptest.NewRecorder()
			v1.class(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_classGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name              string
		withExistingClass bool
		wantResponse      classGetResponse
		wantStatusCode    int
		wantErr           string
	}{
		{
			"request with class in database",
			true,
			classGetResponse{
				newSuccessResponse(),
				model.Class{
					Code:      uuid.NewString(),
					Year:      rand.Int31(),
					Semester:  uuid.NewString(),
					Programme: uuid.NewString(),
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with class not in database",
			false,
			classGetResponse{},
			http.StatusNotFound,
			"the requested class does not exist",
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

			if tt.withExistingClass {
				createdClass := tests.StubClass(t, ctx, v1.db, database.CreateClassParams{
					Code:      tt.wantResponse.Class.Code,
					Year:      tt.wantResponse.Class.Year,
					Semester:  tt.wantResponse.Class.Semester,
					Programme: tt.wantResponse.Class.Programme,
				})

				tt.wantResponse.Class.ID = createdClass.ID
				tt.wantResponse.Class.CreatedAt = createdClass.CreatedAt
				tt.wantResponse.Class.UpdatedAt = createdClass.CreatedAt
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%d", userUrl, tt.wantResponse.Class.ID), nil)
			resp := v1.classGet(req, tt.wantResponse.Class.ID)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGetResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}

func TestAPIServerV1_classPatch(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name               string
		withRequest        classPatchRequest
		withExistingClass  bool
		withUpdateConflict bool
		wantResponse       classPatchResponse
		wantNoChange       bool
		wantStatusCode     int
		wantErr            string
	}{
		{
			"request with field changes",
			classPatchRequest{
				database.UpdateClassParams{
					Programme: to.Ptr("CSC Full-time"),
				},
			},
			true,
			false,
			classPatchResponse{
				newSuccessResponse(),
				classPatchClassResponseFields{
					Code:      uuid.NewString(),
					Year:      rand.Int31(),
					Semester:  uuid.NewString(),
					Programme: "CSC Full-time",
				},
			},
			false,
			http.StatusOK,
			"",
		},
		{
			"request with no field changes",
			classPatchRequest{
				database.UpdateClassParams{},
			},
			true,
			false,
			classPatchResponse{
				newSuccessResponse(),
				classPatchClassResponseFields{
					Code:     uuid.NewString(),
					Year:     rand.Int31(),
					Semester: "2",
				},
			},
			true,
			http.StatusOK,
			"",
		},
		{
			"request updating non-existent class",
			classPatchRequest{
				database.UpdateClassParams{},
			},
			false,
			false,
			classPatchResponse{},
			false,
			http.StatusNotFound,
			"class to update does not exist",
		},
		{
			"request with update conflict",
			classPatchRequest{
				database.UpdateClassParams{
					Code:     to.Ptr(uuid.NewString()),
					Year:     to.Ptr(rand.Int31()),
					Semester: to.Ptr(uuid.NewString()),
				},
			},
			true,
			true,
			classPatchResponse{},
			false,
			http.StatusConflict,
			"class with same code, year, and semester already exists",
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

			var classId int64
			switch {
			case tt.withUpdateConflict:
				// Create the class to update.
				updateClass := tests.StubClass(t, ctx, v1.db, database.CreateClassParams{
					Code:     uuid.NewString(),
					Year:     rand.Int31(),
					Semester: uuid.NewString(),
				})
				classId = updateClass.ID

				// Also create the class to conflict with.
				_ = tests.StubClass(t, ctx, v1.db, database.CreateClassParams{
					Code:     *tt.withRequest.Class.Code,
					Year:     *tt.withRequest.Class.Year,
					Semester: *tt.withRequest.Class.Semester,
				})
			case tt.withExistingClass:
				createdClass := tests.StubClass(t, ctx, v1.db, database.CreateClassParams{
					Code:     tt.wantResponse.Class.Code,
					Year:     tt.wantResponse.Class.Year,
					Semester: tt.wantResponse.Class.Semester,
				})

				classId = createdClass.ID
				tt.wantResponse.Class.ID = createdClass.ID
				tt.wantResponse.Class.UpdatedAt = createdClass.CreatedAt
			default:
				classId = rand.Int63()
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", classUrl, classId), bytes.NewReader(reqBodyBytes))
			resp := v1.classPatch(req, classId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classPatchResponse)
				a.True(ok)

				if !tt.wantNoChange {
					tt.wantResponse.Class.UpdatedAt = actualResp.Class.UpdatedAt
				}

				a.Equal(tt.wantResponse, actualResp)

				// Check that successive updates do not change anything.
				req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", classUrl, classId), bytes.NewReader(reqBodyBytes))
				successiveResp := v1.classPatch(req, classId).(classPatchResponse)
				a.Equal(actualResp, successiveResp)
			}
		})
	}
}

func TestAPIServerV1_classDelete(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                     string
		withExistingClass        bool
		withForeignKeyDependency bool
		wantResponse             classDeleteResponse
		wantStatusCode           int
		wantErr                  string
	}{
		{
			"request with existing class",
			true,
			false,
			classDeleteResponse{newSuccessResponse()},
			http.StatusOK,
			"",
		},
		{
			"request with non-existent class",
			false,
			false,
			classDeleteResponse{},
			http.StatusNotFound,
			"class to delete does not exist",
		},
		{
			"request with class foreign key dependency",
			true,
			true,
			classDeleteResponse{},
			http.StatusConflict,
			"class to delete is still referenced",
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

			var classId int64
			switch {
			case tt.withForeignKeyDependency:
				createdClassGroup := tests.StubClassGroup(t, ctx, v1.db, uuid.NewString(), model.ClassType_Lab)
				classId = createdClassGroup.ClassID
			case tt.withExistingClass:
				createdClass := tests.StubClass(t, ctx, v1.db, database.CreateClassParams{
					Code:      uuid.NewString(),
					Year:      rand.Int31(),
					Semester:  uuid.NewString(),
					Programme: uuid.NewString(),
				})
				classId = createdClass.ID
			default:
				classId = rand.Int63()
			}

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s%d", classUrl, classId), nil)
			resp := v1.classDelete(req, classId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classDeleteResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
