package v1

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) coordinatingClassDashboard(w http.ResponseWriter, r *http.Request, classId int64) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.coordinatingClassDashboardGet(r, classId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type coordinatingClassDashboardGetResponse struct {
	response
	Data database.CoordinatingClassReportData `json:"data"`
}

func (v *APIServerV1) coordinatingClassDashboardGet(r *http.Request, classId int64) apiResponse {
	txDb, tx, err := v.db.AsTx(r.Context(), &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not start database transaction")
	}
	defer func() {
		_ = tx.Rollback()
	}()

	data, err := txDb.GetCoordinatingClassReportData(r.Context(), classId)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusBadRequest, "not allowed to generate report for this class")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not get coordinating class data")
	}

	if err = tx.Commit(); err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not commit database transaction")
	}

	return coordinatingClassDashboardGetResponse{
		newSuccessResponse(),
		data,
	}
}
