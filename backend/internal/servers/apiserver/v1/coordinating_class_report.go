package v1

import (
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/darylhjd/oams/backend/pkg/to"
	"github.com/go-pdf/fpdf"
)

func (v *APIServerV1) coordinatingClassReport(w http.ResponseWriter, r *http.Request, classId int64) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		// Special case for file download, cannot use v.writeResponse helper.
		if err := v.coordinatingClassReportGet(w, r, classId); err != nil {
			resp = *err
		} else {
			return
		}
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

func (v *APIServerV1) coordinatingClassReportGet(w http.ResponseWriter, _ *http.Request, classId int64) *errorResponse {
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{
		"filename": fmt.Sprintf("class_%d_%s_report.pdf", classId, time.Now().Format("2006-01-02_150405")),
	}))
	w.Header().Set("Content-Type", "application/octet-stream")

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	if err := pdf.Output(w); err != nil {
		return to.Ptr(newErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return nil
}
