package v1

import (
	"bytes"
	"context"
	"encoding/json"
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
						StartTime: pgtype.Timestamp{Time: time.UnixMicro(1).UTC(), Valid: true},
						EndTime:   pgtype.Timestamp{Time: time.UnixMicro(2).UTC(), Valid: true},
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

func TestAPIServerV1_classGroupSessionsPost(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withRequest                   classGroupSessionsPostRequest
		withExistingClassGroupSession bool
		withExistingClassGroup        bool
		wantResponse                  classGroupSessionsPostResponse
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with no existing class group session",
			classGroupSessionsPostRequest{
				database.CreateClassGroupSessionParams{
					StartTime: pgtype.Timestamp{Time: time.UnixMicro(1).UTC(), Valid: true},
					EndTime:   pgtype.Timestamp{Time: time.UnixMicro(2).UTC(), Valid: true},
					Venue:     "NEW_CLASS+22",
				},
			},
			false,
			true,
			classGroupSessionsPostResponse{
				newSuccessResponse(),
				database.CreateClassGroupSessionRow{
					Venue: "NEW_CLASS+22",
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with existing class group session",
			classGroupSessionsPostRequest{
				database.CreateClassGroupSessionParams{
					StartTime: pgtype.Timestamp{Time: time.UnixMicro(1).UTC(), Valid: true},
					EndTime:   pgtype.Timestamp{Time: time.UnixMicro(2).UTC(), Valid: true},
					Venue:     "EXISTING_CLASS+22",
				},
			},
			true,
			true,
			classGroupSessionsPostResponse{},
			http.StatusConflict,
			"class group session with same class_group_id and start_time already exists",
		},
		{
			"request with non-existent class group dependency",
			classGroupSessionsPostRequest{
				database.CreateClassGroupSessionParams{
					StartTime: pgtype.Timestamp{Time: time.UnixMicro(1).UTC(), Valid: true},
					EndTime:   pgtype.Timestamp{Time: time.UnixMicro(2).UTC(), Valid: true},
					Venue:     "FAIL_INSERT_CLASS+22",
				},
			},
			false,
			false,
			classGroupSessionsPostResponse{},
			http.StatusBadRequest,
			"class_group_id is not valid",
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
			case tt.withExistingClassGroupSession:
				createdSession := tests.StubClassGroupSession(
					t, ctx, v1.db.Q,
					tt.withRequest.ClassGroupSession.StartTime,
					tt.withRequest.ClassGroupSession.EndTime,
					tt.withRequest.ClassGroupSession.Venue,
				)
				tt.withRequest.ClassGroupSession.ClassGroupID = createdSession.ClassGroupID
			case tt.withExistingClassGroup:
				createdGroup := tests.StubClassGroup(
					t, ctx, v1.db.Q,
					uuid.NewString(),
					database.ClassTypeLAB,
				)
				tt.withRequest.ClassGroupSession.ClassGroupID = createdGroup.ID
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPost, classGroupSessionsUrl, bytes.NewReader(reqBodyBytes))
			resp := v1.classGroupSessionsPost(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupSessionsPostResponse)
				a.True(ok)

				tt.wantResponse.ClassGroupSession.ID = actualResp.ClassGroupSession.ID
				tt.wantResponse.ClassGroupSession.ClassGroupID = actualResp.ClassGroupSession.ClassGroupID
				tt.wantResponse.ClassGroupSession.StartTime = actualResp.ClassGroupSession.StartTime
				tt.wantResponse.ClassGroupSession.EndTime = actualResp.ClassGroupSession.EndTime
				tt.wantResponse.ClassGroupSession.CreatedAt = actualResp.ClassGroupSession.CreatedAt
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
