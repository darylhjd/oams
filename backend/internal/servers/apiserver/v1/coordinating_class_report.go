package v1

import (
	"database/sql"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/darylhjd/oams/backend/internal/servers/apiserver/common"
	"github.com/darylhjd/oams/backend/pkg/to"
	"github.com/go-jet/jet/v2/qrm"
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

func (v *APIServerV1) coordinatingClassReportGet(w http.ResponseWriter, r *http.Request, classId int64) *errorResponse {
	txDb, tx, err := v.db.AsTx(r.Context(), &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		v.logInternalServerError(r, err)
		return to.Ptr(newErrorResponse(http.StatusInternalServerError, "could not start database transaction"))
	}
	defer func() {
		_ = tx.Rollback()
	}()

	data, err := txDb.GetCoordinatingClassReportData(r.Context(), classId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return to.Ptr(newErrorResponse(http.StatusBadRequest, "not allowed to generate report for this class"))
		}

		v.logInternalServerError(r, err)
		return to.Ptr(newErrorResponse(http.StatusInternalServerError, "could not get coordinating class data"))
	}

	if err = tx.Commit(); err != nil {
		v.logInternalServerError(r, err)
		return to.Ptr(newErrorResponse(http.StatusInternalServerError, "could not commit database transaction"))
	}

	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{
		"filename": fmt.Sprintf("class_%d_%s_report.pdf", classId, time.Now().Format("2006-01-02_150405")),
	}))
	w.Header().Set("Content-Type", "application/octet-stream")

	pdf := common.GenerateClassReport(data)
	if err = pdf.Output(w); err != nil {
		return to.Ptr(newErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return nil
}
