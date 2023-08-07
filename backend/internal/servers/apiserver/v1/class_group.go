package v1

import (
	"errors"
	"fmt"
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
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupGet(r, groupId)
	case http.MethodPatch:
		resp = v.classGroupPatch(r, groupId)
	case http.MethodDelete:
		resp = v.classGroupDelete(r, groupId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
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

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group get database action")
	}

	return classGroupGetResponse{
		newSuccessResponse(),
		group,
	}
}

type classGroupPatchRequest struct {
	ClassGroup classGroupPatchClassGroupRequestFields `json:"class_group"`
}

type classGroupPatchClassGroupRequestFields struct {
	ClassID   *int64              `json:"class_id"`
	Name      *string             `json:"name"`
	ClassType *database.ClassType `json:"class_type"`
}

func (r classGroupPatchRequest) updateClassGroupParams(classGroupId int64) database.UpdateClassGroupParams {
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

type classGroupPatchResponse struct {
	response
	ClassGroup database.UpdateClassGroupRow `json:"class_group"`
}

func (v *APIServerV1) classGroupPatch(r *http.Request, id int64) apiResponse {
	var req classGroupPatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	group, err := v.db.Q.UpdateClassGroup(r.Context(), req.updateClassGroupParams(id))
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return newErrorResponse(http.StatusNotFound, "class group to update does not exist")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusBadRequest, "class_id does not exist")
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusConflict, "class group with same class_id, name, and class_type already exists")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process class group patch database action")
		}
	}

	return classGroupPatchResponse{
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
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return newErrorResponse(http.StatusNotFound, "class group to delete does not exist")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusConflict, "class group to delete is still referenced")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process class group delete database action")
		}
	}

	return classGroupDeleteResponse{newSuccessResponse()}
}
