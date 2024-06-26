package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_sessionEnrollmentGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withExistingSessionEnrollment bool
		wantResponse                  sessionEnrollmentGetResponse
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with session enrollment in database",
			true,
			sessionEnrollmentGetResponse{
				newSuccessResponse(),
				model.SessionEnrollment{
					Attended: true,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with session enrollment not in database",
			false,
			sessionEnrollmentGetResponse{},
			http.StatusNotFound,
			"the requested session enrollment does not exist",
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

			if tt.withExistingSessionEnrollment {
				createdEnrollment := tests.StubSessionEnrollment(
					t, ctx, v1.db,
					tt.wantResponse.SessionEnrollment.Attended,
				)

				tt.wantResponse.SessionEnrollment.ID = createdEnrollment.ID
				tt.wantResponse.SessionEnrollment.SessionID = createdEnrollment.SessionID
				tt.wantResponse.SessionEnrollment.UserID = createdEnrollment.UserID
				tt.wantResponse.SessionEnrollment.CreatedAt = createdEnrollment.CreatedAt
				tt.wantResponse.SessionEnrollment.UpdatedAt = createdEnrollment.CreatedAt
			}

			req := httpRequestWithAuthContext(
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%d", sessionEnrollmentUrl, tt.wantResponse.SessionEnrollment.ID), nil),
				tests.StubAuthContext(),
			)
			resp := v1.sessionEnrollmentGet(req, tt.wantResponse.SessionEnrollment.ID)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(sessionEnrollmentGetResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
