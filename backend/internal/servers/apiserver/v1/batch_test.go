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
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/common"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_batch(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name           string
		withMethod     string
		wantStatusCode int
	}{
		{
			"with POST method",
			http.MethodPost,
			http.StatusUnsupportedMediaType,
		},
		{
			"with PUT method",
			http.MethodPut,
			http.StatusBadRequest,
		},
		{
			"with GET method",
			http.MethodGet,
			http.StatusMethodNotAllowed,
		},
		{
			"with PATCH method",
			http.MethodPatch,
			http.StatusMethodNotAllowed,
		},
		{
			"with DELETE method",
			http.MethodDelete,
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
			v1.batch(rr, req)

			a.Equal(tt.wantStatusCode, rr.Code)
		})
	}
}

func TestAPIServerV1_batchPost(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name           string
		body           func() (io.Reader, string, error)
		wantResponse   batchPostResponse
		wantStatusCode int
		wantErr        string
	}{
		{
			"request with file upload",
			func() (io.Reader, string, error) {
				file := "../common/batch_file_well_formatted.xlsx"

				var b bytes.Buffer
				w := multipart.NewWriter(&b)

				if err := w.WriteField(multipartFormBatchWeekIdent, "2"); err != nil {
					return nil, "", err
				}

				f, err := os.Open(file)
				if err != nil {
					return nil, "", err
				}

				ww, err := w.CreateFormFile(multipartFormBatchFileIdent, file)
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

				return &b, w.FormDataContentType(), w.Close()
			},
			batchPostResponse{},
			http.StatusAccepted,
			"",
		},
		{
			"request with non-file content-type",
			func() (io.Reader, string, error) {
				return nil, "application/json", nil
			},
			batchPostResponse{},
			http.StatusUnsupportedMediaType,
			"a multipart request body is required",
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

			body, contentType, err := tt.body()
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPost, batchUrl, body)
			req.Header.Set("Content-Type", contentType)
			resp := v1.batchPost(req)
			a.Equal(tt.wantStatusCode, resp.Code())
		})
	}
}

func TestAPIServerV1_batchPut(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name           string
		body           func() (io.Reader, string, error)
		wantResponse   batchPutResponse
		wantStatusCode int
	}{
		{
			"acceptable request",
			func() (io.Reader, string, error) {
				now := time.Now()

				body := batchPutRequest{
					[]common.BatchData{
						{
							Class: database.UpsertClassParams{
								Code:      "SC1015",
								Year:      2022,
								Semester:  "2",
								Programme: "CSC  Full-Time",
								Au:        3,
							},
							ClassGroups: []common.ClassGroupData{
								{
									database.UpsertClassGroupParams{
										Name:      "A21",
										ClassType: model.ClassType_Lab,
									},
									[]database.UpsertClassGroupSessionParams{
										{
											StartTime: now,
											EndTime:   now.Add(2 * time.Hour),
											Venue:     "",
										},
									},
									[]database.UpsertUserParams{
										{"CHUL6789", "CHUA LI TING"},
										{"YAPW9087", "YAP WEN LI"},
									},
								},
								{
									database.UpsertClassGroupParams{
										Name:      "A26",
										ClassType: model.ClassType_Lab,
									},
									[]database.UpsertClassGroupSessionParams{
										{
											StartTime: now,
											EndTime:   now.Add(2 * time.Hour),
											Venue:     "",
										},
									},
									[]database.UpsertUserParams{
										{"BENST129", "BENJAMIN SANTOS"},
										{"YAPW9088", "YAP WEI LING"},
									},
								},
								{
									database.UpsertClassGroupParams{
										Name:      "A32",
										ClassType: model.ClassType_Lab,
									},
									[]database.UpsertClassGroupSessionParams{
										{
											StartTime: now,
											EndTime:   now.Add(2 * time.Hour),
											Venue:     "",
										},
									},
									[]database.UpsertUserParams{
										{"PATELAR14", "ARJUN PATEL"},
										{"YAPX9087", "YAP XIN TING"},
									},
								},
							},
						},
					},
				}

				b, err := json.Marshal(body)
				if err != nil {
					return nil, "", err
				}

				return bytes.NewReader(b), "application/json", nil
			},
			batchPutResponse{
				newSuccessResponse(),
				[]int64{1},
			},
			http.StatusOK,
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

			body, contentType, err := tt.body()
			a.Nil(err)

			req := httptest.NewRequest(http.MethodPost, batchUrl, body)
			req.Header.Set("Content-Type", contentType)
			resp := v1.batchPut(req)
			a.Equal(tt.wantStatusCode, resp.Code())

			a.Equal(tt.wantResponse, resp)
		})
	}
}
