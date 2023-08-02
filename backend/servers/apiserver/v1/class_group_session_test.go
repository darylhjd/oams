package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func TestAPIServerV1_classGroupSession(t *testing.T) {
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
			http.StatusNotImplemented,
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

			req := httptest.NewRequest(tt.withMethod, fmt.Sprintf("%s%d", classGroupSessionUrl, 1), nil)
			rr := httptest.NewRecorder()
			v1.classGroupSession(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

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
				database.ClassGroupSession{
					StartTime: pgtype.Timestamp{Time: time.UnixMicro(1).UTC(), Valid: true},
					EndTime:   pgtype.Timestamp{Time: time.UnixMicro(2).UTC(), Valid: true},
					Venue:     "EXISTING+46",
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
					t, ctx, v1.db.Q,
					tt.wantResponse.ClassGroupSession.StartTime,
					tt.wantResponse.ClassGroupSession.EndTime,
					tt.wantResponse.ClassGroupSession.Venue,
				)

				tt.wantResponse.ClassGroupSession.ID = createdSession.ID
				tt.wantResponse.ClassGroupSession.ClassGroupID = createdSession.ClassGroupID
				tt.wantResponse.ClassGroupSession.StartTime = createdSession.StartTime
				tt.wantResponse.ClassGroupSession.EndTime = createdSession.EndTime
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

func TestAPIServerV1_classGroupSessionPatch(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withRequest                   classGroupSessionPatchRequest
		withExistingClassGroupSession bool
		wantResponse                  classGroupSessionPatchResponse
		wantNoChange                  bool
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with all fields set",
			classGroupSessionPatchRequest{
				classGroupSessionPatchClassGroupSessionRequestFields{
					ptr(int64(1)),
					ptr(int64(99999999999)),
					ptr(int64(9999999999999)),
					ptr("NEW_VENUE+99"),
				},
			},
			true,
			classGroupSessionPatchResponse{
				newSuccessResponse(),
				database.UpdateClassGroupSessionRow{
					ClassGroupID: 1,
					StartTime:    pgtype.Timestamp{Time: time.UnixMicro(99999999999).UTC(), Valid: true},
					EndTime:      pgtype.Timestamp{Time: time.UnixMicro(9999999999999).UTC(), Valid: true},
					Venue:        "NEW_VENUE+99",
				},
			},
			false,
			http.StatusOK,
			"",
		},
		{
			"request with optional fields not set",
			classGroupSessionPatchRequest{
				classGroupSessionPatchClassGroupSessionRequestFields{},
			},
			true,
			classGroupSessionPatchResponse{
				newSuccessResponse(),
				database.UpdateClassGroupSessionRow{
					StartTime: pgtype.Timestamp{Time: time.UnixMicro(99999999999).UTC(), Valid: true},
					EndTime:   pgtype.Timestamp{Time: time.UnixMicro(9999999999999).UTC(), Valid: true},
					Venue:     "EXISTING_VENUE+99",
				},
			},
			true,
			http.StatusOK,
			"",
		},
		{
			"request updating non-existent class group",
			classGroupSessionPatchRequest{
				classGroupSessionPatchClassGroupSessionRequestFields{},
			},
			false,
			classGroupSessionPatchResponse{
				ClassGroupSession: database.UpdateClassGroupSessionRow{
					ID: 6666,
				},
			},
			false,
			http.StatusNotFound,
			"class group session to update does not exist",
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
					t, ctx, v1.db.Q,
					tt.wantResponse.ClassGroupSession.StartTime,
					tt.wantResponse.ClassGroupSession.EndTime,
					tt.wantResponse.ClassGroupSession.Venue,
				)
				tt.wantResponse.ClassGroupSession.ID = createdSession.ID
				tt.wantResponse.ClassGroupSession.ClassGroupID = createdSession.ClassGroupID
				tt.wantResponse.ClassGroupSession.UpdatedAt = createdSession.CreatedAt
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			sessionId := tt.wantResponse.ClassGroupSession.ID
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", classGroupSessionUrl, sessionId), bytes.NewReader(reqBodyBytes))
			resp := v1.classGroupSessionPatch(req, sessionId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupSessionPatchResponse)
				a.True(ok)

				if !tt.wantNoChange {
					tt.wantResponse.ClassGroupSession.UpdatedAt = actualResp.ClassGroupSession.UpdatedAt
				}

				a.Equal(tt.wantResponse, actualResp)

				// Check that successive updates do not change anything.
				req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", classGroupSessionUrl, sessionId), bytes.NewReader(reqBodyBytes))
				successiveResp := v1.classGroupSessionPatch(req, sessionId).(classGroupSessionPatchResponse)
				a.Equal(actualResp, successiveResp)
			}
		})
	}
}
