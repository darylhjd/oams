package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darylhjd/oams/backend/internal/tests"
)

func TestAPIServerV1_classesCreate(t *testing.T) {
	a := assert.New(t)
	directory := "classes_create"

	files, err := tests.ClassesCreateFiles.ReadDir(directory)
	a.Nil(err)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for _, file := range files {
		f, err := tests.ClassesCreateFiles.Open(filepath.Join(directory, file.Name()))
		a.Nil(err)

		ww, err := w.CreateFormFile(multipartFormFileIdent, file.Name())
		a.Nil(err)

		_, err = io.Copy(ww, f)
		a.Nil(err)

		a.Nil(f.Close())
	}
	a.Nil(w.Close())

	v1 := newTestAPIServerV1(t)
	req := httptest.NewRequest(http.MethodPost, classesUrl, &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rr := httptest.NewRecorder()
	v1.classesCreate(rr, req)

	a.Equal(http.StatusOK, rr.Code)
	expectedBody, err := json.Marshal(classesCreateResponse{
		"file1.txt": fileResult{Success: true},
		"file2.txt": fileResult{Success: true},
	})
	a.Equal(string(expectedBody), rr.Body.String())
}
