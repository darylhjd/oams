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
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
				[]model.SessionEnrollment{
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
				[]model.SessionEnrollment{},
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
					createdEnrollment := tests.StubSessionEnrollment(t, ctx, v1.db, enrollment.Attended)
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

func TestAPIServerV1_sessionEnrollmentsPost(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withRequest                   sessionEnrollmentsPostRequest
		withExistingSessionEnrollment bool
		withExistingClassGroupSession bool
		withExistingUser              bool
		wantResponse                  sessionEnrollmentsPostResponse
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with no existing session enrollment",
			sessionEnrollmentsPostRequest{
				database.CreateSessionEnrollmentParams{
					Attended: true,
				},
			},
			false,
			true,
			true,
			sessionEnrollmentsPostResponse{
				newSuccessResponse(),
				sessionEnrollmentsPostSessionEnrollmentFields{
					Attended: true,
				},
			},
			http.StatusOK,
			"",
		},
		{
			"request with existing session enrollment",
			sessionEnrollmentsPostRequest{
				database.CreateSessionEnrollmentParams{
					Attended: true,
				},
			},
			true,
			true,
			true,
			sessionEnrollmentsPostResponse{},
			http.StatusConflict,
			"session enrollment with same session_id and user_id already exists",
		},
		{
			"request with non existent class group session",
			sessionEnrollmentsPostRequest{
				database.CreateSessionEnrollmentParams{
					Attended: true,
				},
			},
			false,
			false,
			true,
			sessionEnrollmentsPostResponse{},
			http.StatusBadRequest,
			"session_id and/or user_id does not exist",
		},
		{
			"request with non existent user",
			sessionEnrollmentsPostRequest{
				database.CreateSessionEnrollmentParams{
					Attended: true,
				},
			},
			false,
			true,
			false,
			sessionEnrollmentsPostResponse{},
			http.StatusBadRequest,
			"session_id and/or user_id does not exist",
		},
		{
			"request with non existent class group session and user",
			sessionEnrollmentsPostRequest{
				database.CreateSessionEnrollmentParams{
					Attended: true,
				},
			},
			false,
			false,
			false,
			sessionEnrollmentsPostResponse{},
			http.StatusBadRequest,
			"session_id and/or user_id does not exist",
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
			case tt.withExistingSessionEnrollment:
				createdEnrollment := tests.StubSessionEnrollment(t, ctx, v1.db, tt.withRequest.SessionEnrollment.Attended)
				tt.withRequest.SessionEnrollment.SessionID = createdEnrollment.SessionID
				tt.withRequest.SessionEnrollment.UserID = createdEnrollment.UserID
			default:
				if tt.withExistingClassGroupSession {
					createdSession := tests.StubClassGroupSession(
						t, ctx, v1.db.Q,
						pgtype.Timestamptz{Time: time.UnixMicro(1), Valid: true},
						pgtype.Timestamptz{Time: time.UnixMicro(2), Valid: true},
						"VENUE+00",
					)
					tt.withRequest.SessionEnrollment.SessionID = createdSession.ID
				}

				if tt.withExistingUser {
					createdUser := tests.StubUser(t, ctx, v1.db, uuid.NewString(), model.UserRole_Student)
					tt.withRequest.SessionEnrollment.UserID = createdUser.ID
				}
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPost, sessionEnrollmentsUrl, bytes.NewReader(reqBodyBytes))
			resp := v1.sessionEnrollmentsPost(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(sessionEnrollmentsPostResponse)
				a.True(ok)

				tt.wantResponse.SessionEnrollment.ID = actualResp.SessionEnrollment.ID
				tt.wantResponse.SessionEnrollment.SessionID = actualResp.SessionEnrollment.SessionID
				tt.wantResponse.SessionEnrollment.UserID = actualResp.SessionEnrollment.UserID
				tt.wantResponse.SessionEnrollment.CreatedAt = actualResp.SessionEnrollment.CreatedAt
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
