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
	"github.com/darylhjd/oams/backend/internal/servers/apiserver/common"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
		name         string
		body         func() (io.Reader, string, error)
		wantResponse batchPutResponse
	}{
		{
			name: "with json body",
			body: func() (io.Reader, string, error) {
				now := time.Now()

				body := batchPutRequest{
					[]common.BatchData{
						{
							Course: database.UpsertClassesParams{
								Code:      "SC1015",
								Year:      2022,
								Semester:  "2",
								Programme: "CSC  Full-Time",
								Au:        3,
							},
							ClassGroups: []common.ClassGroupData{
								{
									database.UpsertClassGroupsParams{
										Name:      "A21",
										ClassType: database.ClassTypeLAB,
									},
									[]common.SessionData{
										{
											UpsertClassGroupSessionsParams: database.UpsertClassGroupSessionsParams{
												StartTime: pgtype.Timestamp{
													Time:  now,
													Valid: true,
												},
												EndTime: pgtype.Timestamp{
													Time:  now.Add(2 * time.Hour),
													Valid: true,
												},
												Venue: "",
											},
										},
									},
									[]database.UpsertUsersParams{
										{"CHUL6789", "CHUA LI TING", "", database.UserRoleSTUDENT},
										{"YAPW9087", "YAP WEN LI", "", database.UserRoleSTUDENT},
									},
								},
								{
									database.UpsertClassGroupsParams{
										Name:      "A26",
										ClassType: database.ClassTypeLAB,
									},
									[]common.SessionData{
										{
											UpsertClassGroupSessionsParams: database.UpsertClassGroupSessionsParams{
												StartTime: pgtype.Timestamp{
													Time:  now,
													Valid: true,
												},
												EndTime: pgtype.Timestamp{
													Time:  now.Add(2 * time.Hour),
													Valid: true,
												},
												Venue: "",
											},
										},
									},
									[]database.UpsertUsersParams{
										{"BENST129", "BENJAMIN SANTOS", "", database.UserRoleSTUDENT},
										{"YAPW9088", "YAP WEI LING", "", database.UserRoleSTUDENT},
									},
								},
								{
									database.UpsertClassGroupsParams{
										Name:      "A32",
										ClassType: database.ClassTypeLAB,
									},
									[]common.SessionData{
										{
											UpsertClassGroupSessionsParams: database.UpsertClassGroupSessionsParams{
												StartTime: pgtype.Timestamp{
													Time:  now,
													Valid: true,
												},
												EndTime: pgtype.Timestamp{
													Time:  now.Add(2 * time.Hour),
													Valid: true,
												},
												Venue: "",
											},
										},
									},
									[]database.UpsertUsersParams{
										{"PATELAR14", "ARJUN PATEL", "", database.UserRoleSTUDENT},
										{"YAPX9087", "YAP XIN TING", "", database.UserRoleSTUDENT},
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
			wantResponse: batchPutResponse{
				response:           newSuccessResponse(),
				Classes:            1,
				ClassGroups:        3,
				ClassGroupSessions: 3,
				Students:           6,
				SessionEnrollments: 6,
			},
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

			a.Equal(tt.wantResponse, resp)
		})
	}
}
