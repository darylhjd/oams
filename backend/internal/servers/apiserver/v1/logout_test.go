package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/middleware"
	"github.com/darylhjd/oams/backend/internal/oauth2"
)

func TestAPIServerV1_logOut(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name            string
		withAuthContext any
		wantResponse    apiResponse
	}{
		{
			"request with account in context",
			tests.NewMockAuthContext(),
			logoutResponse{newSuccessResponse()},
		},
		{
			"request with wrong account type in context",
			time.Time{},
			newErrorResponse(http.StatusInternalServerError, middleware.ErrUnexpectedAuthContextType.Error()),
		},
		{
			"request with no account in context",
			nil,
			newErrorResponse(http.StatusInternalServerError, "sign-out called but there is no session user"),
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

			req := httptest.NewRequest(http.MethodGet, logoutUrl, nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.AuthContextKey, tt.withAuthContext))
			rr := httptest.NewRecorder()
			v1.logout(rr, req)

			expectedBytes, err := json.Marshal(tt.wantResponse)
			a.Nil(err)
			a.Equal(string(expectedBytes), rr.Body.String())

			if tt.wantResponse.Code() != http.StatusOK {
				return
			}

			// Check for session deletion cookie.
			for _, cookie := range rr.Result().Cookies() {
				if cookie.Name == oauth2.SessionCookieIdent {
					return
				}
			}
			a.FailNow("could not detect expected session deletion cookie")
		})
	}
}
