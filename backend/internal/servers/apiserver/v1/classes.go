package v1

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/darylhjd/oams/backend/internal/database"
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

	v.writeResponse(w, classesUrl, resp)
}

type classesGetResponse struct {
	response
	Classes []database.Class `json:"classes"`
}

func (v *APIServerV1) classesGet(r *http.Request) apiResponse {
	classes, err := v.db.Q.ListClasses(r.Context())
	if err != nil {
		return newErrorResponse(http.StatusInternalServerError, "could not process classes get database action")
	}

	resp := classesGetResponse{
		newSuccessResponse(),
		make([]database.Class, 0, len(classes)),
	}

	resp.Classes = append(resp.Classes, classes...)
	return resp
}

type classesPostRequest struct {
	Class database.CreateClassParams `json:"class"`
}

type classesPostResponse struct {
	response
	Class database.CreateClassRow `json:"class"`
}

func (v *APIServerV1) classesPost(r *http.Request) apiResponse {
	var (
		b   bytes.Buffer
		req classesPostRequest
	)

	if _, err := b.ReadFrom(r.Body); err != nil {
		return newErrorResponse(http.StatusInternalServerError, err.Error())
	}

	if err := json.Unmarshal(b.Bytes(), &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, "could not parse request body")
	}

	class, err := v.db.Q.CreateClass(r.Context(), req.Class)
	if err != nil {
		if database.ErrSQLState(err, database.SQLStateDuplicateKeyOrIndex) {
			return newErrorResponse(http.StatusConflict, "class with same code, year, and semester already exists")
		}

		return newErrorResponse(http.StatusInternalServerError, "could not process classes post database action")
	}

	return classesPostResponse{
		newSuccessResponse(),
		class,
	}
}
