package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/go-jet/jet/v2/qrm"
)

func (v *APIServerV1) classGroupManager(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	managerId, err := strconv.ParseInt(r.PathValue("managerId"), 10, 64)
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
	Manager model.ClassGroupManager `json:"manager"`
}

func (v *APIServerV1) classGroupManagerGet(r *http.Request, managerId int64) apiResponse {
	manager, err := v.db.GetClassGroupManager(r.Context(), managerId)
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
	ManagingRole model.ManagingRole `json:"managing_role"`
}

type classGroupManagerPatchResponse struct {
	response
	Manager model.ClassGroupManager `json:"manager"`
}

func (v *APIServerV1) classGroupManagerPatch(r *http.Request, managerId int64) apiResponse {
	var req classGroupManagerPatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	manager, err := v.db.UpdateClassGroupManager(r.Context(), database.UpdateClassGroupManagerParams{
		ID:   managerId,
		Role: req.ManagingRole,
	})
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusUnauthorized, "not allowed to update managing role")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not update class group manager role")
	}

	return classGroupManagerPatchResponse{
		newSuccessResponse(),
		manager,
	}
}

type classGroupManagerDeleteResponse struct {
	response
}

func (v *APIServerV1) classGroupManagerDelete(r *http.Request, managerId int64) apiResponse {
	if err := v.db.DeleteClassGroupManager(r.Context(), managerId); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusUnauthorized, "not allowed to delete manager")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not delete class group manager")
	}

	return classGroupManagerDeleteResponse{
		newSuccessResponse(),
	}
}
