package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_classGroupSessions(t *testing.T) {
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

			req := httptest.NewRequest(tt.withMethod, classGroupSessionsUrl, nil)
			rr := httptest.NewRecorder()
			v1.classGroupSessions(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_classGroupSessionsGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withExistingClassGroupSession bool
		wantResponse                  classGroupSessionsGetResponse
	}{
		{
			"request with class group session in database",
			true,
			classGroupSessionsGetResponse{
				newSuccessResponse(),
				[]database.ClassGroupSession{
					{
						StartTime: pgtype.Timestamp{Time: time.Now(), Valid: true},
						EndTime:   pgtype.Timestamp{Time: time.Now(), Valid: true},
						Venue:     "CLASS+22",
					},
				},
			},
		},
		{
			"request with no class group session in database",
			false,
			classGroupSessionsGetResponse{
				newSuccessResponse(),
				[]database.ClassGroupSession{},
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

			if tt.withExistingClassGroupSession {
				for idx, session := range tt.wantResponse.ClassGroupSessions {
					createdSession := tests.StubClassGroupSession(t, ctx, v1.db.Q, session.StartTime, session.EndTime, session.Venue)
					sessionPtr := &tt.wantResponse.ClassGroupSessions[idx]
					sessionPtr.ID = createdSession.ID
					sessionPtr.ClassGroupID = createdSession.ClassGroupID
					sessionPtr.StartTime, sessionPtr.EndTime = createdSession.StartTime, createdSession.EndTime
					sessionPtr.CreatedAt, sessionPtr.UpdatedAt = createdSession.CreatedAt, createdSession.CreatedAt
				}
			}

			req := httptest.NewRequest(http.MethodGet, classGroupSessionsUrl, nil)
			actualResp, ok := v1.classGroupSessionsGet(req).(classGroupSessionsGetResponse)
			a.True(ok)
			a.Equal(tt.wantResponse, actualResp)
		})
	}
}
