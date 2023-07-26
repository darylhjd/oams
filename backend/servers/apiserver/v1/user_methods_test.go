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
	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func Test_userGet(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name            string
		withAuthContext any
		wantResponse    userGetResponse
		wantErr         string
	}{
		{
			"request with account in context",
			tests.NewMockAuthContext(),
			userGetResponse{
				response: newSuccessResponse(),
				SessionUser: &database.User{
					ID: tests.MockAuthenticatorIDTokenName,
					Email: pgtype.Text{
						String: tests.MockAuthenticatorAccountPreferredUsername,
						Valid:  true,
					},
				},
				Users: []database.User{},
			},
			"",
		},
		{
			"request with no account in context",
			nil,
			userGetResponse{
				response: newSuccessResponse(),
				Users:    []database.User{},
			},
			"",
		},
		{
			"request with wrong auth context type",
			time.Time{},
			userGetResponse{},
			middleware.ErrUnexpectedAuthContextType.Error(),
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

			tests.StubAuthContextUser(t, ctx, v1.db.Q)

			req := httptest.NewRequest(http.MethodGet, userUrl, nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.AuthContextKey, tt.withAuthContext))
			actualResp := v1.userGet(req)

			switch {
			case tt.wantErr != "":
				err, ok := actualResp.(errorResponse)
				a.True(ok)
				a.Contains(err.Error, tt.wantErr)
				return
			case tt.withAuthContext != nil:
				resp, ok := actualResp.(userGetResponse)
				a.True(ok)
				a.Equal(tt.wantResponse.SessionUser.ID, resp.SessionUser.ID)
				tt.wantResponse.SessionUser = resp.SessionUser
				fallthrough
			default:
				a.Equal(tt.wantResponse, actualResp)
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
