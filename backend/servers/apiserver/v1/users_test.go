package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/oauth2"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/middleware"
)

func TestAPIServerV1_users(t *testing.T) {
	tts := []struct {
		name            string
		withAuthContext any
		wantResponse    apiResponse
	}{
		{
			"request with account in context",
			middleware.AuthContext{
				AuthResult: oauth2.NewMockAzureAuthenticator().MockAuthResult(),
			},
			usersResponse{
				response: newSuccessfulResponse(),
				Users:    []database.Student{},
			},
		},
		{
			"request with no account in context",
			nil,
			usersResponse{
				response: newSuccessfulResponse(),
				Users:    []database.Student{},
			},
		},
		{
			"request with wrong account type in context",
			time.Time{},
			newErrorResponse(http.StatusInternalServerError, middleware.ErrUnexpectedAuthContextType.Error()),
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			id := uuid.NewString()

			v1 := newTestAPIServerV1(t, id)
			defer tests.TearDown(t, v1.db, id)

			err := v1.db.Q.UpsertStudents(context.Background(), []database.UpsertStudentsParams{
				tests.MockUpsertStudentsParams(),
			}).Close()
			a.Nil(err)

			req := httptest.NewRequest(http.MethodGet, usersUrl, nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.AuthContextKey, tt.withAuthContext))
			rr := httptest.NewRecorder()
			v1.users(rr, req)

			a.Equal(tt.wantResponse.Code(), rr.Code)

			var expectedResp apiResponse
			if rr.Code != http.StatusOK {
				expectedResp = tt.wantResponse
			} else {
				actualResp := usersResponse{}
				a.Nil(json.Unmarshal(rr.Body.Bytes(), &actualResp))

				var sessionUser *database.Student
				if tt.withAuthContext != nil {
					t.Log(actualResp.SessionUser)
					authResult := tt.withAuthContext.(middleware.AuthContext).AuthResult
					sessionUser = &database.Student{
						ID: authResult.IDToken.Name,
						Email: pgtype.Text{
							String: authResult.Account.PreferredUsername,
							Valid:  true,
						},
						CreatedAt: actualResp.SessionUser.CreatedAt,
						UpdatedAt: actualResp.SessionUser.UpdatedAt,
					}
				}

				expected := tt.wantResponse.(usersResponse)
				expected.SessionUser = sessionUser
				expectedResp = expected
			}

			bytes, err := json.Marshal(expectedResp)
			a.Nil(err)
			a.Equal(string(bytes), rr.Body.String())
		})
	}
}
