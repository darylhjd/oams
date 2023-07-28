package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (v *APIServerV1) class(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	classId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, classUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, classUrl, newErrorResponse(http.StatusUnprocessableEntity, "invalid class id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classGet(r, classId)
	case http.MethodPut:
		resp = v.classPut(r, classId)
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, classUrl, resp)
}

type classGetResponse struct {
	response
	Class database.Class `json:"class"`
}

func (v *APIServerV1) classGet(r *http.Request, id int64) apiResponse {
	class, err := v.db.Q.GetClass(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process class get database action")
	}

	return classGetResponse{
		newSuccessResponse(),
		class,
	}
}

type classPutRequest struct {
	Class classPutClassRequestFields `json:"class"`
}

type classPutClassRequestFields struct {
	Code      *string `json:"code"`
	Year      *int32  `json:"year"`
	Semester  *string `json:"semester"`
	Programme *string `json:"programme"`
	Au        *int16  `json:"au"`
}

type classPutResponse struct {
	response
	Class database.UpdateClassRow `json:"class"`
}

func (r *classPutRequest) updateClassParams(classId int64) database.UpdateClassParams {
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

func (v *APIServerV1) classPut(r *http.Request, id int64) apiResponse {
	var (
		b   bytes.Buffer
		req classPutRequest
	)

	if _, err := b.ReadFrom(r.Body); err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if err := json.Unmarshal(b.Bytes(), &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, "could not parse request body")
	}

	class, err := v.db.Q.UpdateClass(r.Context(), req.updateClassParams(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "class to update does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process class put database action")
	}

	return classPutResponse{
		newSuccessResponse(),
		class,
	}
}
