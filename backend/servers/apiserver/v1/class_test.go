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
				database.Class{
					Code:     "CZ3454",
					Year:     2023,
					Semester: "1",
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
				createdClass := tests.StubClass(
					t, ctx, v1.db.Q,
					tt.wantResponse.Class.Code,
					tt.wantResponse.Class.Year,
					tt.wantResponse.Class.Semester,
				)

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
		name              string
		withRequest       classPatchRequest
		withExistingClass bool
		wantResponse      classPatchResponse
		wantNoChange      bool
		wantStatusCode    int
		wantErr           string
	}{
		{
			"request with all fields set",
			classPatchRequest{
				classPatchClassRequestFields{
					ptr("CZ9999"),
					ptr(int32(1999)),
					ptr("1"),
					ptr("CSC Full-time"),
					ptr(int16(3)),
				},
			},
			true,
			classPatchResponse{
				newSuccessResponse(),
				database.UpdateClassRow{
					Code:      "CZ9999",
					Year:      1999,
					Semester:  "1",
					Programme: "CSC Full-time",
					Au:        3,
				},
			},
			false,
			http.StatusOK,
			"",
		},
		{
			"request with optional fields not set",
			classPatchRequest{
				classPatchClassRequestFields{},
			},
			true,
			classPatchResponse{
				newSuccessResponse(),
				database.UpdateClassRow{
					Code:     "EXISTING123",
					Year:     2023,
					Semester: "1",
				},
			},
			true,
			http.StatusOK,
			"",
		},
		{
			"request updating non-existent class",
			classPatchRequest{
				classPatchClassRequestFields{},
			},
			false,
			classPatchResponse{
				Class: database.UpdateClassRow{
					ID: 6666,
				},
			},
			false,
			http.StatusNotFound,
			"class to update does not exist",
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
				createdClass := tests.StubClass(
					t, ctx, v1.db.Q,
					tt.wantResponse.Class.Code,
					tt.wantResponse.Class.Year,
					tt.wantResponse.Class.Semester,
				)
				tt.wantResponse.Class.ID = createdClass.ID
				tt.wantResponse.Class.UpdatedAt = createdClass.CreatedAt
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			classId := tt.wantResponse.Class.ID
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
		name              string
		withExistingClass bool
		wantResponse      classDeleteResponse
		wantStatusCode    int
		wantErr           string
	}{
		{
			"request with existing class",
			true,
			classDeleteResponse{newSuccessResponse()},
			http.StatusOK,
			"",
		},
		{
			"request with non-existent class",
			false,
			classDeleteResponse{},
			http.StatusNotFound,
			"class to delete does not exist",
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

			var classId int64 = 6666 // Choose a random ID that does not exist.
			if tt.withExistingClass {
				createdClass := tests.StubClass(t, ctx, v1.db.Q, "RANDOM_CODE", 9999, "22")
				classId = createdClass.ID
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
