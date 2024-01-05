package v1

import (
	"fmt"
	"net/http"
	"strings"
)

func (v *APIServerV1) signature(w http.ResponseWriter, r *http.Request) {
	var resp apiResponse

	userId := strings.TrimPrefix(r.URL.Path, signatureUrl)
	switch r.Method {
	case http.MethodPut:
		resp = v.signaturePut(r, userId)
	default:
		resp = newErrorResponse(http.StatusMethodNotAllowed, "")
	}

	v.writeResponse(w, r, resp)
}

type signaturePutRequest struct {
	Signature string `json:"signature"`
}

type signaturePutResponse struct {
	response
}

func (v *APIServerV1) signaturePut(r *http.Request, userId string) apiResponse {
	var req signaturePutRequest
	if err := v.parseRequestBody(r.Body, &req); err != nil {
		return newErrorResponse(http.StatusBadRequest, fmt.Sprintf("could not parse user request body: %s", err))
	}

	if err := v.db.UpdateUserSignature(r.Context(), userId, req.Signature); err != nil {
		v.logInternalServerError(r, err)
		return newErrorResponse(http.StatusInternalServerError, "could not update user signature")
	}

	return signaturePutResponse{newSuccessResponse()}
}
