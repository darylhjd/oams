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

func (v *APIServerV1) classGroupManager(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	managerId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, classGroupManagerUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class group manager id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupManagerGet(r, managerId)
	case http.MethodPatch:
		resp = v.classGroupManagerPatch(r, managerId)
	case http.MethodDelete:
		resp = v.classGroupManagerDelete(r, managerId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGroupManagerGetResponse struct {
	response
	ClassGroupManager model.ClassGroupManager `json:"class_group_manager"`
}

func (v *APIServerV1) classGroupManagerGet(r *http.Request, id int64) apiResponse {
	manager, err := v.db.GetClassGroupManager(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class group manager does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group manager get database action")
	}

	return classGroupManagerGetResponse{
		newSuccessResponse(),
		manager,
	}
}

type classGroupManagerPatchRequest struct {
	ClassGroupManager database.UpdateClassGroupManagerParams `json:"class_group_manager"`
}

type classGroupManagerPatchResponse struct {
	response
	ClassGroupManager classGroupManagerPatchClassGroupManagerResponseFields `json:"class_group_manager"`
}

type classGroupManagerPatchClassGroupManagerResponseFields struct {
	ID           int64              `json:"id"`
	UserID       string             `json:"user_id"`
	ClassGroupID int64              `json:"class_group_id"`
	ManagingRole model.ManagingRole `json:"managing_role"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

func (v *APIServerV1) classGroupManagerPatch(r *http.Request, id int64) apiResponse {
	var req classGroupManagerPatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	manager, err := v.db.UpdateClassGroupManager(r.Context(), id, req.ClassGroupManager)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "class group manager to update does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group manager patch database action")
	}

	return classGroupManagerPatchResponse{
		newSuccessResponse(),
		classGroupManagerPatchClassGroupManagerResponseFields{
			manager.ID,
			manager.UserID,
			manager.ClassGroupID,
			manager.ManagingRole,
			manager.UpdatedAt,
		},
	}
}

type classGroupManagerDeleteResponse struct {
	response
}

func (v *APIServerV1) classGroupManagerDelete(r *http.Request, id int64) apiResponse {
	_, err := v.db.DeleteClassGroupManager(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "class group manager to delete does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group manager delete database action")
	}

	return classGroupManagerDeleteResponse{newSuccessResponse()}
}
