package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/tests"
	"github.com/darylhjd/oams/backend/servers/apiserver/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestAPIServerV1_classesCreate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tts := []struct {
		name         string
		body         func() (io.Reader, string, error)
		wantResponse classesCreateResponse
	}{
		{
			"with file upload",
			func() (io.Reader, string, error) {
				file := "class_lab_test.xlsx"

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
			classesCreateResponse{
				newSuccessfulResponse(),
				1,
				3,
				3,
				58,
				66,
			},
		},
		{
			name: "with json body",
			body: func() (io.Reader, string, error) {
				now := time.Now()

				body := classesCreateRequest{
					[]common.ClassCreationData{
						{
							Course: database.UpsertCoursesParams{
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
									[]database.UpsertStudentsParams{
										{"CHUL6789", "CHUA LI TING", pgtype.Text{}},
										{"YAPW9087", "YAP WEN LI", pgtype.Text{}},
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
									[]database.UpsertStudentsParams{
										{"BENST129", "BENJAMIN SANTOS", pgtype.Text{}},
										{"YAPW9088", "YAP WEI LING", pgtype.Text{}},
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
									[]database.UpsertStudentsParams{
										{"PATELAR14", "ARJUN PATEL", pgtype.Text{}},
										{"YAPX9087", "YAP XIN TING", pgtype.Text{}},
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
			wantResponse: classesCreateResponse{
				response:           newSuccessfulResponse(),
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

			body, contentType, err := tt.body()
			if err != nil {
				t.Fatal(err)
			}

			v1 := newTestAPIServerV1(t, id)
			defer tests.TearDown(t, v1.db, id)

			req := httptest.NewRequest(http.MethodPost, classesUrl, body)
			req.Header.Set("Content-Type", contentType)
			rr := httptest.NewRecorder()
			v1.classesCreate(rr, req)

			b, err := json.Marshal(tt.wantResponse)
			a.Nil(err)
			a.Equal(string(b), rr.Body.String())

			// Check correct number of inputs in database.
			courses, err := v1.db.Q.ListCourses(ctx)
			a.Nil(err)
			a.Equal(tt.wantResponse.Classes, len(courses))

			classGroups, err := v1.db.Q.ListClassGroups(ctx)
			a.Nil(err)
			a.Equal(tt.wantResponse.ClassGroups, len(classGroups))

			classGroupSessions, err := v1.db.Q.ListClassGroupSessions(ctx)
			a.Nil(err)
			a.Equal(tt.wantResponse.ClassGroupSessions, len(classGroupSessions))

			students, err := v1.db.Q.ListStudents(ctx)
			a.Nil(err)
			a.Equal(tt.wantResponse.Students, len(students))

			sessionEnrollments, err := v1.db.Q.ListSessionEnrollments(ctx)
			a.Nil(err)
			a.Equal(tt.wantResponse.SessionEnrollments, len(sessionEnrollments))
		})
	}
}
