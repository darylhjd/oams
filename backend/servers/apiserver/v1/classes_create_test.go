package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_classesCreate(t *testing.T) {
	tests := []struct {
		name         string
		body         func() (io.Reader, string, error)
		expectedBody any
	}{
		{
			"with files",
			func() (io.Reader, string, error) {
				files := []string{
					"class_lab_test.xlsx",
					"class_lec_test.xlsx",
				}

				var b bytes.Buffer
				w := multipart.NewWriter(&b)
				for _, file := range files {
					f, err := os.Open(file)
					if err != nil {
						return nil, "", err
					}

					ww, err := w.CreateFormFile(multipartFormFileIdent, file)
					if err != nil {
						return nil, "", err
					}

					_, err = io.Copy(ww, f)
					if err != nil {
						return nil, "", err
					}

					if err = f.Close(); err != nil {
						return nil, "", err
					}
				}

				return &b, w.FormDataContentType(), w.Close()
			},
			classesCreateResponse{
				"class_lab_test.xlsx": fileProcessingResult{Success: true},
				"class_lec_test.xlsx": fileProcessingResult{Success: true},
			},
		},
	}

	a := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, contentType, err := tt.body()
			if err != nil {
				t.Fatal(err)
			}

			v1 := newTestAPIServerV1(t)
			req := httptest.NewRequest(http.MethodPost, classesUrl, body)
			req.Header.Set("Content-Type", contentType)
			rr := httptest.NewRecorder()
			v1.classesCreate(rr, req)

			a.Equal(http.StatusOK, rr.Code)
			expectedBody, err := json.Marshal(tt.expectedBody)
			if err != nil {
				t.Fatal(err)
			}

			a.Equal(string(expectedBody), rr.Body.String())
		})
	}
}
