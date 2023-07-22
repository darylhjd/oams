package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"

	"github.com/darylhjd/oams/backend/servers/apiserver/common"
	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/pkg/goroutines"
)

const (
	maxParseMemory         = 32 << 20
	maxGoRoutines          = 10
	multipartFormFileIdent = "attachments"
)

type classesCreateRequest struct {
	Classes []common.ClassCreationData `json:"classes"`
}

// classesCreateResponse is a data type detailing the result of the classes create endpoint.
type classesCreateResponse []common.ClassCreationData

// classesCreate is the handler for a request to create classes.
func (v *APIServerV1) classesCreate(w http.ResponseWriter, r *http.Request) {
	var (
		resp classesCreateResponse
		err  error
	)
	switch contentType := r.Header.Get("Content-Type"); {
	case strings.HasPrefix(contentType, "multipart"):
		resp, err = v.processClassCreationFiles(r)
	case contentType == "application/json":
		resp, err = v.processClassCreationJSON(r)
	default:
		v.l.Debug(fmt.Sprintf("%s - received classes create request with unacceptable content-type", namespace),
			zap.String("content-type", contentType))
		http.Error(w, "unacceptable content-type for classes creation request", http.StatusUnsupportedMediaType)
		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("url", classesUrl),
			zap.Error(err))
	}
}

// processClassCreationFiles processes a request to create classes via file uploads.
func (v *APIServerV1) processClassCreationFiles(r *http.Request) (classesCreateResponse, error) {
	if err := r.ParseMultipartForm(maxParseMemory); err != nil {
		return nil, err
	}

	limiter := goroutines.NewLimiter(maxGoRoutines)

	saveRes := sync.Map{}
	for _, header := range r.MultipartForm.File[multipartFormFileIdent] {
		header := header // Required for go routine to point to different file for each loop.
		limiter.Do(func() {
			creationData, err := v.processClassCreationFile(header)
			saveRes.Store(&creationData, err)
		})
	}

	limiter.Wait()

	var (
		toProcess []common.ClassCreationData
		err       error
	)
	saveRes.Range(func(key, value any) bool {
		data, ok := key.(*common.ClassCreationData)
		if !ok {
			err = errors.New("type assertion failed when processing class creation data")
			return false
		}

		if value != nil {
			err = value.(error)
			return false
		}

		toProcess = append(toProcess, *data)
		return true
	})

	if err != nil {
		return nil, err
	}

	processingResp, err := v.processClasses(toProcess)
	if err != nil {
		return nil, err
	}

	return processingResp, err
}

// processClassCreationFile processes a file to create a new class.
func (v *APIServerV1) processClassCreationFile(fileHeader *multipart.FileHeader) (common.ClassCreationData, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return common.ClassCreationData{}, err
	}

	v.l.Debug(fmt.Sprintf("%s - processing class creation file", namespace),
		zap.String("filename", fileHeader.Filename))

	creationData, err := common.ParseClassCreationFile(fileHeader.Filename, file)
	if err != nil {
		return common.ClassCreationData{}, fmt.Errorf("%s - error parsing class creation file %s: %w", namespace, fileHeader.Filename, err)
	}

	return *creationData, file.Close()
}

// processClassCreationJSON processes a request to create classes via JSON body.
func (v *APIServerV1) processClassCreationJSON(r *http.Request) (classesCreateResponse, error) {
	var b bytes.Buffer
	if _, err := b.ReadFrom(r.Body); err != nil {
		return nil, err
	}

	var request classesCreateRequest
	if err := json.Unmarshal(b.Bytes(), &request); err != nil {
		return nil, err
	}

	return v.processClasses(request.Classes)
}

// processClasses sequentially processes each class creation data provided.
func (v *APIServerV1) processClasses(classes []common.ClassCreationData) (classesCreateResponse, error) {
	var resp classesCreateResponse
	for idx, _ := range classes {
		class := classes[idx]
		resp = append(resp, class)
		if !(&class).Validate() { // Defensively check for validity of creation data.
			continue
		}

		// TODO: Implement database action for inserting classes.
	}

	return resp, nil
}
