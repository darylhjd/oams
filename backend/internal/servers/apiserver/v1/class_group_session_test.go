package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
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
			http.StatusNotFound,
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
					StartTime: pgtype.Timestamptz{Time: time.UnixMicro(1), Valid: true},
					EndTime:   pgtype.Timestamptz{Time: time.UnixMicro(2), Valid: true},
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
		withUpdateConflict            bool
		withExistingUpdateClassGroup  bool
		wantResponse                  classGroupSessionPatchResponse
		wantNoChange                  bool
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with field changes",
			classGroupSessionPatchRequest{
				classGroupSessionPatchClassGroupSessionRequestFields{
					ptr(int64(1)),
					ptr(int64(99999999999)),
					ptr(int64(9999999999999)),
					ptr("NEW_VENUE+99"),
				},
			},
			true,
			false,
			true,
			classGroupSessionPatchResponse{
				newSuccessResponse(),
				database.UpdateClassGroupSessionRow{
					ClassGroupID: 1,
					StartTime:    pgtype.Timestamptz{Time: time.UnixMicro(99999999999), Valid: true},
					EndTime:      pgtype.Timestamptz{Time: time.UnixMicro(9999999999999), Valid: true},
					Venue:        "NEW_VENUE+99",
				},
			},
			false,
			http.StatusOK,
			"",
		},
		{
			"request with no field changes",
			classGroupSessionPatchRequest{
				classGroupSessionPatchClassGroupSessionRequestFields{},
			},
			true,
			false,
			true,
			classGroupSessionPatchResponse{
				newSuccessResponse(),
				database.UpdateClassGroupSessionRow{
					StartTime: pgtype.Timestamptz{Time: time.UnixMicro(99999999999), Valid: true},
					EndTime:   pgtype.Timestamptz{Time: time.UnixMicro(9999999999999), Valid: true},
					Venue:     "EXISTING_VENUE+99",
				},
			},
			true,
			http.StatusOK,
			"",
		},
		{
			"request updating non-existent class group session",
			classGroupSessionPatchRequest{
				classGroupSessionPatchClassGroupSessionRequestFields{},
			},
			false,
			false,
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
		{
			"request with update conflict",
			classGroupSessionPatchRequest{
				classGroupSessionPatchClassGroupSessionRequestFields{
					StartTime: ptr(int64(2)),
					EndTime:   ptr(int64(3)),
				},
			},
			true,
			true,
			true,
			classGroupSessionPatchResponse{},
			false,
			http.StatusConflict,
			"class group session with same class_group_id and start_time already exists",
		},
		{
			"request with non-existent class group dependency",
			classGroupSessionPatchRequest{},
			true,
			false,
			false,
			classGroupSessionPatchResponse{},
			false,
			http.StatusBadRequest,
			"class_group_id does not exist",
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

			var sessionId int64
			switch {
			case tt.withUpdateConflict:
				// Create session to update.
				updateClassGroupSession := tests.StubClassGroupSession(
					t, ctx, v1.db.Q,
					pgtype.Timestamptz{Time: time.UnixMicro(1), Valid: true},
					pgtype.Timestamptz{Time: time.UnixMicro(2), Valid: true},
					uuid.NewString(),
				)
				sessionId = updateClassGroupSession.ID

				// Also create session to conflict with.
				_ = tests.StubClassGroupSessionWithClassGroupID(
					t, ctx, v1.db.Q,
					updateClassGroupSession.ClassGroupID,
					pgtype.Timestamptz{Time: time.UnixMicro(*tt.withRequest.ClassGroupSession.StartTime), Valid: true},
					pgtype.Timestamptz{Time: time.UnixMicro(*tt.withRequest.ClassGroupSession.EndTime), Valid: true},
					uuid.NewString(),
				)
			case tt.withExistingClassGroupSession && !tt.withExistingUpdateClassGroup:
				createdSession := tests.StubClassGroupSession(
					t, ctx, v1.db.Q,
					pgtype.Timestamptz{Time: time.UnixMicro(1), Valid: true},
					pgtype.Timestamptz{Time: time.UnixMicro(2), Valid: true},
					uuid.NewString(),
				)

				sessionId = createdSession.ID
				tt.withRequest.ClassGroupSession.ClassGroupID = ptr(createdSession.ClassGroupID + 1)
			case tt.withExistingClassGroupSession:
				createdSession := tests.StubClassGroupSession(
					t, ctx, v1.db.Q,
					tt.wantResponse.ClassGroupSession.StartTime,
					tt.wantResponse.ClassGroupSession.EndTime,
					tt.wantResponse.ClassGroupSession.Venue,
				)

				sessionId = createdSession.ID
				tt.wantResponse.ClassGroupSession.ID = createdSession.ID
				tt.wantResponse.ClassGroupSession.ClassGroupID = createdSession.ClassGroupID
				tt.wantResponse.ClassGroupSession.UpdatedAt = createdSession.CreatedAt
			default:
				sessionId = rand.Int63()
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

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

func TestAPIServerV1_classGroupSessionDelete(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withExistingClassGroupSession bool
		withForeignKeyDependency      bool
		wantResponse                  classGroupSessionDeleteResponse
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with existing class group session",
			true,
			false,
			classGroupSessionDeleteResponse{newSuccessResponse()},
			http.StatusOK,
			"",
		},
		{
			"request with non-existent class group",
			false,
			false,
			classGroupSessionDeleteResponse{},
			http.StatusNotFound,
			"class group session to delete does not exist",
		},
		{
			"request with class group session foreign key dependency",
			true,
			true,
			classGroupSessionDeleteResponse{},
			http.StatusConflict,
			"class group session to delete is still referenced",
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

			var sessionId int64
			switch {
			case tt.withForeignKeyDependency:
				createdSessionEnrollment := tests.StubSessionEnrollment(t, ctx, v1.db, true)
				sessionId = createdSessionEnrollment.SessionID
			case tt.withExistingClassGroupSession:
				createdSession := tests.StubClassGroupSession(
					t, ctx, v1.db.Q,
					pgtype.Timestamptz{Time: time.UnixMicro(1), Valid: true},
					pgtype.Timestamptz{Time: time.UnixMicro(2), Valid: true},
					uuid.NewString(),
				)
				sessionId = createdSession.ID
			default:
				sessionId = rand.Int63()
			}

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s%d", classGroupSessionUrl, sessionId), nil)
			resp := v1.classGroupSessionDelete(req, sessionId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(classGroupSessionDeleteResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
