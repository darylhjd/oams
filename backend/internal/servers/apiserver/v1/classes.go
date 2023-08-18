package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/darylhjd/oams/backend/internal/database"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/model"
	"github.com/darylhjd/oams/backend/internal/database/gen/oams/public/table"
)

func (v *APIServerV1) classes(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	switch r.Method {
	case http.MethodGet:
		resp = v.classesGet(r)
	case http.MethodPost:
		resp = v.classesPost(r)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type classesGetResponse struct {
	response
	Classes []model.Class `json:"classes"`
}

func (v *APIServerV1) classesGet(r *http.Request) apiResponse {
	params, err := v.decodeListQueryParameters(r.URL.Query(), table.Classes.AllColumns)
	if err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}

	classes, err := v.db.ListClasses(r.Context(), params)
	if err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process classes get database action")
	}

	resp := classesGetResponse{
		newSuccessResponse(),
		make([]model.Class, 0, len(classes)),
	}

	resp.Classes = append(resp.Classes, classes...)
	return resp
}

type classesPostRequest struct {
	Class database.CreateClassParams `json:"class"`
}

type classesPostResponse struct {
	response
	Class classesPostClassResponseFields `json:"class"`
}

type classesPostClassResponseFields struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Year      int32     `json:"year"`
	Semester  string    `json:"semester"`
	Programme string    `json:"programme"`
	Au        int16     `json:"au"`
	CreatedAt time.Time `json:"created_at"`
}

func (v *APIServerV1) classesPost(r *http.Request) apiResponse {
	var req classesPostRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse request body: %s", err))
	}

	class, err := v.db.CreateClass(r.Context(), req.Class)
	if err != nil {
		if database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex) {
			return newErrorResponse(http.StatusConflict, "class with same code, year, and semester already exists")
		}

		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not process classes post database action")
	}

	return classesPostResponse{
		newSuccessResponse(),
		classesPostClassResponseFields{
			class.ID,
			class.Code,
			class.Year,
			class.Semester,
			class.Programme,
			class.Au,
			class.CreatedAt,
		},
	}
}
