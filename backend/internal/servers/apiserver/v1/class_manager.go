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

func (v *APIServerV1) classManager(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	managerId, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, classManagerUrl), 10, 64)
	if err != nil {
		v.writeResponse(w, r, newErrorResponse(http.StatusUnprocessableEntity, "invalid class manager id"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		resp = v.classManagerGet(r, managerId)
	case http.MethodPatch:
		resp = v.classManagerPatch(r, managerId)
	case http.MethodDelete:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classManagerGetResponse struct {
	response
	ClassManager model.ClassManager `json:"class_manager"`
}

func (v *APIServerV1) classManagerGet(r *http.Request, id int64) apiResponse {
	manager, err := v.db.GetClassManager(r.Context(), id)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "the requested class manager does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class manager get database action")
	}

	return classManagerGetResponse{
		newSuccessResponse(),
		manager,
	}
}

type classManagerPatchRequest struct {
	ClassManager database.UpdateClassManagerParams `json:"class_manager"`
}

type classManagerPatchResponse struct {
	response
	ClassManager classManagerPatchClassManagerResponseFields `json:"class_manager"`
}

type classManagerPatchClassManagerResponseFields struct {
	ID           int64              `json:"id"`
	UserID       string             `json:"user_id"`
	ClassID      int64              `json:"class_id"`
	ManagingRole model.ManagingRole `json:"managing_role"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

func (v *APIServerV1) classManagerPatch(r *http.Request, id int64) apiResponse {
	var req classManagerPatchRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request bodyL %s", err))
	}

	manager, err := v.db.UpdateClassManager(r.Context(), id, req.ClassManager)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return newErrorResponse(http.StatusNotFound, "class manager to update does not exist")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class manager patch database action")
	}

	return classManagerPatchResponse{
		newSuccessResponse(),
		classManagerPatchClassManagerResponseFields{
			manager.ID,
			manager.UserID,
			manager.ClassID,
			manager.ManagingRole,
			manager.UpdatedAt,
		},
	}
}
