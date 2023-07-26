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
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_user(t *testing.T) {
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
			"with PUT method",
			http.MethodPut,
			http.StatusBadRequest,
		},
		{
			"with DELETE method",
			http.MethodDelete,
			http.StatusNotImplemented,
		},
		{
			"with PATCH method",
			http.MethodPatch,
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

			req := httptest.NewRequest(tt.withMethod, userUrl, nil)
			rr := httptest.NewRecorder()
			v1.user(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_userGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name           string
		withUserId     string
		wantStatusCode int
		wantErr        string
	}{
		{
			"request with user in database",
			uuid.NewString(),
			http.StatusOK,
			"",
		},
		{
			"request with user not in database",
			uuid.NewString(),
			http.StatusNotFound,
			"the requested user does not exist",
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

			if tt.wantErr == "" {
				tests.StubUser(t, ctx, v1.db.Q, tt.withUserId)
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", userUrl, tt.withUserId), nil)
			resp := v1.userGet(req, tt.withUserId)
			a.Equal(tt.wantStatusCode, resp.Code())

			switch {
			case tt.wantErr != "":
				actualResp, ok := resp.(errorResponse)
				a.True(ok)
				a.Contains(actualResp.Error, tt.wantErr)
			default:
				actualResp, ok := resp.(userGetResponse)
				a.True(ok)
				a.Equal(tt.withUserId, actualResp.User.ID)
			}
		})
	}
}

func Test_userPut(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name         string
		withRequest  any
		wantResponse userPutResponse
		wantErr      string
	}{
		{
			"good request",
			userPutRequest{
				[]database.UpsertUsersParams{
					{
						ID:   tests.MockAuthenticatorIDTokenName,
						Name: "",
						Email: pgtype.Text{
							String: tests.MockAuthenticatorAccountPreferredUsername,
							Valid:  true,
						},
						Role: database.UserRoleSTUDENT,
					},
				},
			},
			userPutResponse{
				newSuccessResponse(),
				[]database.User{
					{
						ID:   tests.MockAuthenticatorIDTokenName,
						Name: "",
						Email: pgtype.Text{
							String: tests.MockAuthenticatorAccountPreferredUsername,
							Valid:  true,
						},
						Role: database.UserRoleSTUDENT,
					},
				},
			},
			"",
		},
		{
			"bad request",
			[]string{},
			userPutResponse{},
			"could not parse request body",
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

			reqBodyBytes, err := json.Marshal(tt.withRequest)
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPut, userUrl, bytes.NewReader(reqBodyBytes))
			actualResp := v1.userPut(req)

			switch {
			case tt.wantErr != "":
				err, ok := actualResp.(errorResponse)
				a.True(ok)
				a.Contains(err.Error, tt.wantErr)
			default:
				resp, ok := actualResp.(userPutResponse)
				a.True(ok)
				a.Equal(len(tt.wantResponse.Users), len(resp.Users))
				for idx, respUser := range resp.Users {
					tt.wantResponse.Users[idx].CreatedAt, tt.wantResponse.Users[idx].UpdatedAt = respUser.CreatedAt, respUser.UpdatedAt
				}

				a.Equal(tt.wantResponse, resp)
			}
		})
	}
}
