package v1

import (
	"context"
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
