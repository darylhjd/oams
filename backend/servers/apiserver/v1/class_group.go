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

func (v *APIServerV1) classGroup(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	groupId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, classGroupUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, classGroupUrl, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupGet(r, groupId)
	case http.MethodPut:
		resp = v.classGroupPut(r, groupId)
	case http.MethodDelete:
		resp = v.classGroupDelete(r, groupId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, classGroupUrl, resp)
}

type classGroupGetResponse struct {
	response
	ClassGroup database.ClassGroup `json:"class_group"`
}

func (v *APIServerV1) classGroupGet(r *http.Request, id int64) apiResponse {
	group, err := v.db.Q.GetClassGroup(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class group does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process class group get database action")
	}

	return classGroupGetResponse{
		newSuccessResponse(),
		group,
	}
}

type classGroupPutRequest struct {
	ClassGroup classGroupPutClassGroupRequestFields `json:"class_group"`
}

type classGroupPutClassGroupRequestFields struct {
	ClassID   *int64              `json:"class_id"`
	Name      *string             `json:"name"`
	ClassType *database.ClassType `json:"class_type"`
}

func (r classGroupPutRequest) updateClassGroupParams(classGroupId int64) database.UpdateClassGroupParams {
	params := database.UpdateClassGroupParams{ID: classGroupId}

	if r.ClassGroup.ClassID != nil {
		params.ClassID = pgtype.Int8{Int64: *r.ClassGroup.ClassID, Valid: true}
	}

	if r.ClassGroup.Name != nil {
		params.Name = pgtype.Text{String: *r.ClassGroup.Name, Valid: true}
	}

	if r.ClassGroup.ClassType != nil {
		params.ClassType = database.NullClassType{ClassType: *r.ClassGroup.ClassType, Valid: true}
	}

	return params
}

type classGroupPutResponse struct {
	response
	ClassGroup database.UpdateClassGroupRow `json:"class_group"`
}

func (v *APIServerV1) classGroupPut(r *http.Request, id int64) apiResponse {
	var (
		b   bytes.Buffer
		req classGroupPutRequest
	)

	if _, err := b.ReadFrom(r.Body); err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if err := json.Unmarshal(b.Bytes(), &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, "could not parse request body")
	}

	group, err := v.db.Q.UpdateClassGroup(r.Context(), req.updateClassGroupParams(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "class group to update does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process class group put database action")
	}

	return classGroupPutResponse{
		newSuccessResponse(),
		group,
	}
}

type classGroupDeleteResponse struct {
	response
}

func (v *APIServerV1) classGroupDelete(r *http.Request, id int64) apiResponse {
	_, err := v.db.Q.DeleteClassGroup(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "class group to delete does not exist")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process class group delete database action")
	}

	return classGroupDeleteResponse{newSuccessResponse()}
}
