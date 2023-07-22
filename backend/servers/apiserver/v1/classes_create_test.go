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
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/servers/apiserver/common"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_classesCreate(t *testing.T) {
	tests := []struct {
		name         string
		body         func() (io.Reader, string, error)
		expectedBody classesCreateResponse
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
				{Filename: "class_lab_test.xlsx"},
				{Filename: "class_lec_test.xlsx"},
			},
		},
		{
			name: "with json body",
			body: func() (io.Reader, string, error) {
				body := classesCreateRequest{
					[]common.ClassCreationData{
						{
							"class_lab_test.xlsx",
							time.Date(2023, time.June, 15, 13, 1, 0, 0, time.Local),
							database.CreateCoursesParams{
								Code:      "SC1015",
								Year:      2022,
								Semester:  "2",
								Programme: "CSC  Full-Time",
								Au:        3,
							},
							[]common.ClassGroupData{
								{
									database.CreateClassGroupsParams{
										Name:      "A21",
										ClassType: database.ClassTypeLAB,
									},
									[]database.CreateClassGroupSessionsParams{},
									[]database.CreateStudentsParams{
										{"CHUL6789", "CHUA LI TING", pgtype.Text{}},
										{"YAPW9087", "YAP WEN LI", pgtype.Text{}},
									},
								},
								{
									database.CreateClassGroupsParams{
										Name:      "A26",
										ClassType: database.ClassTypeLAB,
									},
									[]database.CreateClassGroupSessionsParams{},
									[]database.CreateStudentsParams{
										{"BENST129", "BENJAMIN SANTOS", pgtype.Text{}},
										{"YAPW9087", "YAP WEI LING", pgtype.Text{}},
									},
								},
								{
									database.CreateClassGroupsParams{
										Name:      "A32",
										ClassType: database.ClassTypeLAB,
									},
									[]database.CreateClassGroupSessionsParams{},
									[]database.CreateStudentsParams{
										{"PATELAR14", "ARJUN PATEL", pgtype.Text{}},
										{"YAPX9087", "YAP XIN TING", pgtype.Text{}},
									},
								},
							},
							"",
							database.ClassTypeLAB,
						},
					},
				}

				b, err := json.Marshal(body)
				if err != nil {
					return nil, "", err
				}

				return bytes.NewReader(b), "application/json", nil
			},
			expectedBody: classesCreateResponse{
				{Course: database.CreateCoursesParams{}},
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

			var actualResponse classesCreateResponse
			a.Nil(json.Unmarshal(rr.Body.Bytes(), &actualResponse))
			a.Equal(len(tt.expectedBody), len(actualResponse))
		})
	}
}
