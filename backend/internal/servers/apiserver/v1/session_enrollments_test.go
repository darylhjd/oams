package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_sessionEnrollments(t *testing.T) {
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

			req := httptest.NewRequest(tt.withMethod, sessionEnrollmentsUrl, nil)
			rr := httptest.NewRecorder()
			v1.sessionEnrollments(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_sessionEnrollmentsGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withExistingSessionEnrollment bool
		wantResponse                  sessionEnrollmentsGetResponse
	}{
		{
			"request with session enrollment in database",
			true,
			sessionEnrollmentsGetResponse{
				newSuccessResponse(),
				[]database.SessionEnrollment{
					{
						Attended: false,
					},
				},
			},
		},
		{
			"request with no session enrollment in database",
			false,
			sessionEnrollmentsGetResponse{
				newSuccessResponse(),
				[]database.SessionEnrollment{},
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

			if tt.withExistingSessionEnrollment {
				for idx, enrollment := range tt.wantResponse.SessionEnrollments {
					createdEnrollment := tests.StubSessionEnrollment(t, ctx, v1.db.Q, enrollment.Attended)
					sessionPtr := &tt.wantResponse.SessionEnrollments[idx]
					sessionPtr.ID = createdEnrollment.ID
					sessionPtr.SessionID = createdEnrollment.SessionID
					sessionPtr.UserID = createdEnrollment.UserID
					sessionPtr.CreatedAt, sessionPtr.UpdatedAt = createdEnrollment.CreatedAt, createdEnrollment.CreatedAt
				}
			}

			req := httptest.NewRequest(http.MethodGet, sessionEnrollmentsUrl, nil)
			actualResp, ok := v1.sessionEnrollmentsGet(req).(sessionEnrollmentsGetResponse)
			a.True(ok)
			a.Equal(tt.wantResponse, actualResp)
		})
	}
}
