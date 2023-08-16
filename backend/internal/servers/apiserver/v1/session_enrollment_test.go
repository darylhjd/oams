package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_sessionEnrollment(t *testing.T) {
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

			req := httptest.NewRequest(tt.withMethod, fmt.Sprintf("%s%d", sessionEnrollmentUrl, 1), nil)
			rr := httptest.NewRecorder()
			v1.sessionEnrollment(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%d", sessionEnrollmentUrl, tt.wantResponse.SessionEnrollment.ID), nil)
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

func TestAPIServerV1_sessionEnrollmentPatch(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withRequest                   sessionEnrollmentPatchRequest
		withExistingSessionEnrollment bool
		wantResponse                  sessionEnrollmentPatchResponse
		wantNoChange                  bool
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with field changes",
			sessionEnrollmentPatchRequest{
				sessionEnrollmentPatchSessionEnrollmentRequestFields{
					Attended: ptr(true),
				},
			},
			true,
			sessionEnrollmentPatchResponse{
				newSuccessResponse(),
				database.UpdateSessionEnrollmentRow{
					Attended: true,
				},
			},
			false,
			http.StatusOK,
			"",
		},
		{
			"request with no field changes",
			sessionEnrollmentPatchRequest{
				sessionEnrollmentPatchSessionEnrollmentRequestFields{},
			},
			true,
			sessionEnrollmentPatchResponse{
				newSuccessResponse(),
				database.UpdateSessionEnrollmentRow{
					Attended: true,
				},
			},
			true,
			http.StatusOK,
			"",
		},
		{
			"request updating non-existent session enrollment",
			sessionEnrollmentPatchRequest{
				sessionEnrollmentPatchSessionEnrollmentRequestFields{},
			},
			false,
			sessionEnrollmentPatchResponse{
				SessionEnrollment: database.UpdateSessionEnrollmentRow{
					ID: 6666,
				},
			},
			false,
			http.StatusNotFound,
			"session enrollment to update does not exist",
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
				createdEnrollment := tests.StubSessionEnrollment(t, ctx, v1.db, tt.wantResponse.SessionEnrollment.Attended)
				tt.wantResponse.SessionEnrollment.ID = createdEnrollment.ID
				tt.wantResponse.SessionEnrollment.SessionID = createdEnrollment.SessionID
				tt.wantResponse.SessionEnrollment.UserID = createdEnrollment.UserID
				tt.wantResponse.SessionEnrollment.UpdatedAt = pgtype.Timestamptz{Time: createdEnrollment.CreatedAt, Valid: true}
			}

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			enrollmentId := tt.wantResponse.SessionEnrollment.ID
			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", sessionEnrollmentUrl, enrollmentId), bytes.NewReader(reqBodyBytes))
			resp := v1.sessionEnrollmentPatch(req, enrollmentId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(sessionEnrollmentPatchResponse)
				a.True(ok)

				if !tt.wantNoChange {
					tt.wantResponse.SessionEnrollment.UpdatedAt = actualResp.SessionEnrollment.UpdatedAt
				}

				a.Equal(tt.wantResponse, actualResp)

				// Check that successive updates do not change anything.
				req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("%s%d", sessionEnrollmentUrl, enrollmentId), bytes.NewReader(reqBodyBytes))
				successiveResp := v1.sessionEnrollmentPatch(req, enrollmentId).(sessionEnrollmentPatchResponse)
				a.Equal(actualResp, successiveResp)
			}
		})
	}
}

func TestAPIServerV1_sessionEnrollmentDelete(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name                          string
		withExistingSessionEnrollment bool
		wantResponse                  sessionEnrollmentDeleteResponse
		wantStatusCode                int
		wantErr                       string
	}{
		{
			"request with existing session enrollment",
			true,
			sessionEnrollmentDeleteResponse{newSuccessResponse()},
			http.StatusOK,
			"",
		},
		{
			"request with non-existent session enrollment",
			false,
			sessionEnrollmentDeleteResponse{},
			http.StatusNotFound,
			"session enrollment to delete does not exist",
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

			var enrollmentId int64 = 6666 // Choose a random ID that does not exist.
			if tt.withExistingSessionEnrollment {
				createdEnrollment := tests.StubSessionEnrollment(t, ctx, v1.db, false)
				enrollmentId = createdEnrollment.ID
			}

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s%d", sessionEnrollmentUrl, enrollmentId), nil)
			resp := v1.sessionEnrollmentDelete(req, enrollmentId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(sessionEnrollmentDeleteResponse)
				a.True(ok)
				a.Equal(tt.wantResponse, actualResp)
			}
		})
	}
}
