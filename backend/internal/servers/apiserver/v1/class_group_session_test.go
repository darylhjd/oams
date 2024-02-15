package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_classGroupSessionGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withExistingClassGroupSession bool
		wantResponse                  classGroupSessionGetResponse
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with class group session in database",
			true,
			classGroupSessionGetResponse{
				newSuccessResponse(),
				model.ClassGroupSession{
					StartTime: time.UnixMicro(999),
					EndTime:   time.UnixMicro(99999),
					Venue:     uuid.NewString(),
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with class group session not in database",
			false,
			classGroupSessionGetResponse{},
			http.StatusNotFound,
			"the requested class group session does not exist",
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

			if tt.withExistingClassGroupSession {
				createdSession := tests.StubClassGroupSession(
					t, ctx, v1.db,
					tt.wantResponse.ClassGroupSession.StartTime,
					tt.wantResponse.ClassGroupSession.EndTime,
					tt.wantResponse.ClassGroupSession.Venue,
				)

				tt.wantResponse.ClassGroupSession.ID = createdSession.ID
				tt.wantResponse.ClassGroupSession.ClassGroupID = createdSession.ClassGroupID
				tt.wantResponse.ClassGroupSession.CreatedAt = createdSession.CreatedAt
				tt.wantResponse.ClassGroupSession.UpdatedAt = createdSession.CreatedAt
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%d", classGroupSessionUrl, tt.wantResponse.ClassGroupSession.ID), nil)
			resp := v1.classGroupSessionGet(req, tt.wantResponse.ClassGroupSession.ID)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupSessionGetResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
