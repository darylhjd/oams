package v1

import (
	"fmt"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
)

func (v *APIServerV1) classGroups(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupsGet(r)
	case http.MethodPost:
		resp = v.classGroupsPost(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGroupsGetResponse struct {
	response
	ClassGroups []database.ClassGroup `json:"class_groups"`
}

func (v *APIServerV1) classGroupsGet(r *http.Request) apiResponse {
	groups, err := v.db.Q.ListClassGroups(r.Context())
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class groups get database action")
	}

	resp := classGroupsGetResponse{
		newSuccessResponse(),
		make([]database.ClassGroup, 0, len(groups)),
	}

	resp.ClassGroups = append(resp.ClassGroups, groups...)
	return resp
}

type classGroupsPostRequest struct {
	ClassGroup database.CreateClassGroupParams `json:"class_group"`
}

type classGroupsPostResponse struct {
	response
	ClassGroup database.CreateClassGroupRow `json:"class_group"`
}

func (v *APIServerV1) classGroupsPost(r *http.Request) apiResponse {
	var req classGroupsPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	group, err := v.db.Q.CreateClassGroup(r.Context(), req.ClassGroup)
	if err != nil {
		switch {
		case database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex):
			return newErrorResponse(http.StatusConflict, "class group with same class_id, name, and class_type already exists")
		case database.ErrSQLState(err, database.SQLStateForeignKeyViolation):
			return newErrorResponse(http.StatusBadRequest, "class_id does not exist")
		default:
			v.logInternalServerError(r, err)
			return newErrorResponse(http.StatusInternalServerError, "could not process class groups post database action")
		}
	}

	return classGroupsPostResponse{
		newSuccessResponse(),
		group,
	}
}
