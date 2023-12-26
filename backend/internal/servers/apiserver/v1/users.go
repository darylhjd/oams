package v1

import (
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/postgres/public/table"
)

func (v *APIServerV1) users(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.usersGet(r)
	case http.MethodPost:
		resp = newErrorResponse(http.StatusNotImplemented, "")
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type usersGetResponse struct {
	response
	Users []model.User `json:"users"`
}

func (v *APIServerV1) usersGet(r *http.Request) apiResponse {
	params, err := database.DecodeListQueryParams(r.URL.Query(), table.Users.AllColumns)
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}

	users, err := v.db.ListUsers(r.Context(), params)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process users get database action")
	}

	resp := usersGetResponse{
		newSuccessResponse(),
		make([]model.User, 0, len(users)),
	}

	resp.Users = append(resp.Users, users...)
	return resp
}
