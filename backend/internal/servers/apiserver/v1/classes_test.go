package v1

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestAPIServerV1_classes(t *testing.T) {
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

			req := httptest.NewRequest(tt.withMethod, classesUrl, nil)
			rr := httptest.NewRecorder()
			v1.classes(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_classesGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name              string
		withExistingClass bool
		wantResponse      classesGetResponse
	}{
		{
			"request with class in database",
			true,
			classesGetResponse{
				newSuccessResponse(),
				[]model.Class{
					{
						Code:     "CZ1115",
						Year:     2023,
						Semester: "2",
					},
				},
			},
		},
		{
			"request with no class in database",
			false,
			classesGetResponse{
				newSuccessResponse(),
				[]model.Class{},
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

			if tt.withExistingClass {
				for idx, class := range tt.wantResponse.Classes {
					createdClass := tests.StubClass(t, ctx, v1.db, class.Code, class.Year, class.Semester)
					classPtr := &tt.wantResponse.Classes[idx]
					classPtr.ID = createdClass.ID
					classPtr.CreatedAt, classPtr.UpdatedAt = createdClass.CreatedAt, createdClass.CreatedAt
				}
			}

			req := httptest.NewRequest(http.MethodGet, classesUrl, nil)
			actualResp, ok := v1.classesGet(req).(classesGetResponse)
			a.True(ok)
			a.Equal(tt.wantResponse, actualResp)
		})
	}
}

func TestAPIServerV1_classesGetQueryParams(t *testing.T) {
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
				"sort": []string{"au"},
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

			u := url.URL{Path: classesUrl}
			u.RawQuery = tt.query.Encode()

			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			resp := v1.classesGet(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				_, ok := resp.(classesGetResponse)
				a.True(ok)
			}
		})
	}
}

func TestAPIServerV1_classesPost(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name              string
		withRequest       classesPostRequest
		withExistingClass bool
		wantResponse      classesPostResponse
		wantStatusCode    int
		wantErr           string
	}{
		{
			"request with no existing class",
			classesPostRequest{
				database.CreateClassParams{
					Code:     "CZ1115",
					Year:     2023,
					Semester: "2",
				},
			},
			false,
			classesPostResponse{
				newSuccessResponse(),
				classesPostClassResponseFields{
					Code:     "CZ1115",
					Year:     2023,
					Semester: "2",
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with existing class",
			classesPostRequest{
				database.CreateClassParams{
					Code:     "CZ1115",
					Year:     2023,
					Semester: "2",
				},
			},
			true,
			classesPostResponse{},
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

			if tt.withExistingClass {
				_ = tests.StubClass(
					t, ctx, v1.db,
					tt.withRequest.Class.Code,
					tt.withRequest.Class.Year,
					tt.withRequest.Class.Semester,
				)
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPost, classesUrl, bytes.NewReader(reqBodyBytes))
			resp := v1.classesPost(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classesPostResponse)
				a.True(ok)

				tt.wantResponse.Class.ID = actualResp.Class.ID
				tt.wantResponse.Class.CreatedAt = actualResp.Class.CreatedAt
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
