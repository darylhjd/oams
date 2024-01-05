package v1

import (
	"database/sql"
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/darylhjd/oams/backend/internal/servers/apiserver/common"
	"github.com/darylhjd/oams/backend/pkg/to"
)

func (v *APIServerV1) dataExport(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		// Special case for file download, cannot use v.writeResponse helper.
		if err := v.dataExportGet(w, r); err != nil {
			resp = *err
		} else {
			return
		}
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

func (v *APIServerV1) dataExportGet(w http.ResponseWriter, r *http.Request) *errorResponse {
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

	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{
		"filename": fmt.Sprintf("%s_oams_snapshot.zip", time.Now().Format("2006-01-02_150405")),
	}))
	w.Header().Set("Content-Type", "application/octet-stream")

	if err := common.GenerateDataExportArchive(w, r, txDb); err != nil {
		v.logInternalServerError(r, err)
		return to.Ptr(newErrorResponse(http.StatusInternalServerError, "could not generate data export zip"))
	}

	if err = tx.Commit(); err != nil {
		v.logInternalServerError(r, err)
		return to.Ptr(newErrorResponse(http.StatusInternalServerError, "could not commit database transaction"))
	}

	return nil
}
