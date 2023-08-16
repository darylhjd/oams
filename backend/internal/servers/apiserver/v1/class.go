package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/go-jet/jet/v2/qrm"
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
	Class database.UpdateClassParams `json:"class"`
}

type classPatchResponse struct {
	response
	Class classPatchClassResponseFields `json:"class"`
}

type classPatchClassResponseFields struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Year      int32     `json:"year"`
	Semester  string    `json:"semester"`
	Programme string    `json:"programme"`
	Au        int16     `json:"au"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (v *APIServerV1) classPatch(r *http.Request, id int64) apiResponse {
	var req classPatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	class, err := v.db.UpdateClass(r.Context(), id, req.Class)
	if err != nil {
		switch {
		case errors.Is(err, qrm.ErrNoRows):
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
		classPatchClassResponseFields{
			class.ID,
			class.Code,
			class.Year,
			class.Semester,
			class.Programme,
			class.Au,
			class.UpdatedAt,
		},
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
