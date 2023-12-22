package v1

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
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
						Code:      uuid.NewString(),
						Year:      rand.Int31(),
						Semester:  uuid.NewString(),
						Programme: uuid.NewString(),
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
					createdClass := tests.StubClass(t, ctx, v1.db, database.CreateClassParams{
						Code:      class.Code,
						Year:      class.Year,
						Semester:  class.Semester,
						Programme: class.Programme,
					})
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
