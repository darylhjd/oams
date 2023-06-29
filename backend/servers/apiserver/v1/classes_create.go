package v1

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"go.uber.org/zap"

	"github.com/darylhjd/oams/backend/pkg/goroutines"
)

const (
	maxParseMemory         = 32 << 20
	maxGoRoutines          = 10
	multipartFormFileIdent = "attachments"
)

// classesCreateResponse is a data type detailing the result of the classes create endpoint.
type classesCreateResponse map[string]fileResult

// fileResult holds the result of processing a file.
type fileResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// classesCreate is the handler for a request to create classes.
func (v *APIServerV1) classesCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxParseMemory); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	limiter := goroutines.NewLimiter(maxGoRoutines)
	results := classesCreateResponse{}
	for _, header := range r.MultipartForm.File[multipartFormFileIdent] {
		header := header // Required for go routine to point to different file for each loop.
		limiter.Do(func() {
			res := fileResult{Success: true}
			if err := v.processClassCreateFile(header); err != nil {
				res.Success = false
				res.Error = err.Error()
			}

			results[header.Filename] = res
		})
	}
	limiter.Wait()

	bytes, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(bytes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		v.l.Error(fmt.Sprintf("%s - could not write response", namespace),
			zap.String("url", classesUrl),
			zap.Error(err))
	}
}

// processClassCreateFile processes a file to create a new class.
// TODO: Switch this with proper implementation once format of file is provided.
func (v *APIServerV1) processClassCreateFile(fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	v.l.Debug(fmt.Sprintf("%s - processing class create file", namespace),
		zap.String("filename", fileHeader.Filename))
	return file.Close()
}
