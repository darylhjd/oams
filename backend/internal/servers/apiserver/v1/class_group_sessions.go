package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
)

func (v *APIServerV1) classGroupSessions(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.classGroupSessionsGet(r)
	case http.MethodPost:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classGroupSessionsGetResponse struct {
	response
	ClassGroupSessions []model.ClassGroupSession `json:"class_group_sessions"`
}

func (v *APIServerV1) classGroupSessionsGet(r *http.Request) apiResponse {
	params, err := database.DecodeListQueryParams(r.URL.Query(), table.ClassGroupSessions.AllColumns)
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}

	sessions, err := v.db.ListClassGroupSessions(r.Context(), params)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process class group sessions get database action")
	}

	return classGroupSessionsGetResponse{
		newSuccessResponse(),
		append(make([]model.ClassGroupSession, 0, len(sessions)), sessions...),
	}
}
