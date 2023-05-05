package servers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllowMethods(t *testing.T) {
	tests := []struct {
		name           string
		testMethods    []string
		allowedMethods []string
		wantErr        bool
	}{
		{
			"success on allowed methods",
			[]string{
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
				http.MethodConnect,
				http.MethodHead,
				http.MethodOptions,
				http.MethodPatch,
				http.MethodTrace,
				http.MethodPut},
			[]string{
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
				http.MethodConnect,
				http.MethodHead,
				http.MethodOptions,
				http.MethodPatch,
				http.MethodTrace,
				http.MethodPut,
			},
			false,
		},
		{
			"fail on disallowed method",
			[]string{http.MethodPost},
			[]string{http.MethodGet},
			true,
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, method := range tt.testMethods {
				req, err := http.NewRequest(method, "", nil)
				if err != nil {
					t.Fatal(err)
				}

				testHandler := AllowMethods(func(w http.ResponseWriter, r *http.Request) {}, tt.allowedMethods...)
				rr := httptest.NewRecorder()

				testHandler.ServeHTTP(rr, req)

				a.Equal(tt.wantErr, rr.Result().StatusCode == http.StatusMethodNotAllowed)
			}
		})
	}
}