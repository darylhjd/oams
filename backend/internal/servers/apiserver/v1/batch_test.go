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

func TestAPIServerV1_batchPost(t *testing.T) {
	t.Parallel()

	tts := []struct {
		name         string
		body         func() (io.Reader, string, error)
		wantResponse batchPostResponse
	}{
		{
			"with file upload",
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
			batchPostResponse{
				newSuccessResponse(),
				1,
				3,
				4,
				58,
				94,
			},
		},
		{
			name: "with json body",
			body: func() (io.Reader, string, error) {
				now := time.Now()

				body := batchPostRequest{
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
			wantResponse: batchPostResponse{
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
			resp := v1.batchPost(req)

			a.Equal(tt.wantResponse, resp)
		})
	}
}
