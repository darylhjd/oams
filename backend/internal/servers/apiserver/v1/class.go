package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (v *APIServerV1) class(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, classUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classGet(r, classId)
	case http.MethodPatch:
		resp = v.classPatch(r, classId)
	case http.MethodDelete:
		resp = v.classDelete(r, classId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGetResponse struct {
	response
	Class model.Class `json:"class"`
}

func (v *APIServerV1) classGet(r *http.Request, id int64) apiResponse {
	class, err := v.db.GetClass(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class get database action")
	}

	return classGetResponse{
		newSuccessResponse(),
		class,
	}
}

type classPatchRequest struct {
	Class classPatchClassRequestFields `json:"class"`
}

type classPatchClassRequestFields struct {
	Code      *string `json:"code"`
	Year      *int32  `json:"year"`
	Semester  *string `json:"semester"`
	Programme *string `json:"programme"`
	Au        *int16  `json:"au"`
}

func (r classPatchRequest) updateClassParams(classId int64) database.UpdateClassParams {
	params := database.UpdateClassParams{ID: classId}

	if r.Class.Code != nil {
		params.Code = pgtype.Text{String: *r.Class.Code, Valid: true}
	}

	if r.Class.Year != nil {
		params.Year = pgtype.Int4{Int32: *r.Class.Year, Valid: true}
	}

	if r.Class.Semester != nil {
		params.Semester = pgtype.Text{String: *r.Class.Semester, Valid: true}
	}

	if r.Class.Programme != nil {
		params.Programme = pgtype.Text{String: *r.Class.Programme, Valid: true}
	}

	if r.Class.Au != nil {
		params.Au = pgtype.Int2{Int16: *r.Class.Au, Valid: true}
	}

	return params
}

type classPatchResponse struct {
	response
	Class database.UpdateClassRow `json:"class"`
}

func (v *APIServerV1) classPatch(r *http.Request, id int64) apiResponse {
	var req classPatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	class, err := v.db.Q.UpdateClass(r.Context(), req.updateClassParams(id))
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return newErrorResponse(http.StatusNotFound, "class to update does not exist")
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusConflict, "class with same code, year, and semester already exists")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process class patch database action")
		}
	}

	return classPatchResponse{
		newSuccessResponse(),
		class,
	}
}

type classDeleteResponse struct {
	response
}

func (v *APIServerV1) classDelete(r *http.Request, id int64) apiResponse {
	_, err := v.db.DeleteClass(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, qrm.ErrNoRows):
			return newErrorResponse(http.StatusNotFound, "class to delete does not exist")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusConflict, "class to delete is still referenced")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process class delete database action")
		}
	}

	return classDeleteResponse{newSuccessResponse()}
}
