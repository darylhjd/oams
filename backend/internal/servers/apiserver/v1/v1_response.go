package v1

import "net/http"

// apiResponse is an interface defining all responses that the server must fulfill.
type apiResponse interface {
	Code() int
}

// response struct contains fields that all responses from the API must have.
type response struct {
	Result bool `json:"result"`

	// This field is used to set the response status code and does not appear in the response body.
	statusCode int
}

func (r response) Code() int {
	return r.statusCode
}

// newSuccessResponse creates a new response struct with true result and StatusOK status code.
func newSuccessResponse() response {
	return response{true, http.StatusOK}
}

// errorResponse struct contains fields that all error responses from the API must have.
type errorResponse struct {
	response
	Error string `json:"error"`
}

// newErrorResponse creates a new errorResponse. Caller may specify the status code and the error message.
// The client is not expected to read the body if the status code is not StatusOK.
func newErrorResponse(code int, err string) errorResponse {
	return errorResponse{
		response{false, code},
		err,
	}
}
