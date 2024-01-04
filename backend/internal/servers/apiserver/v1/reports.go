package v1

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/darylhjd/oams/backend/pkg/to"
)

func (v *APIServerV1) reports(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		// Special case for file download, cannot use v.writeResponse helper.
		if err := v.reportsGet(w, r); err != nil {
			resp = *err
		} else {
			return
		}
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

func (v *APIServerV1) reportsGet(w http.ResponseWriter, _ *http.Request) *errorResponse {
	var buffer bytes.Buffer

	records := [][]string{
		{"This is a test file, with a comma and a double quote \"", "This is the next field."},
	}

	csvWriter := csv.NewWriter(&buffer)
	if err := csvWriter.WriteAll(records); err != nil {
		return to.Ptr(
			newErrorResponse(http.StatusInternalServerError, fmt.Sprintf("could not create csv file: %s", err)),
		)
	}

	filename := fmt.Sprintf("%s_snapshot.csv", time.Now().Format("2006-01-02_150405"))
	cd := mime.FormatMediaType("attachment", map[string]string{
		"filename": filename,
	})

	w.Header().Set("Content-Disposition", cd)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))
	if _, err := io.Copy(w, &buffer); err != nil {
		return to.Ptr(
			newErrorResponse(http.StatusInternalServerError, fmt.Sprintf("failed to write file: %s", err)),
		)
	}

	return nil
}
