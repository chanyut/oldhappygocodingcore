package core

import (
	"net/http"
)

type APIError struct {
	ErrorID              int    `json:"id"`
	HTTPStatus           int    `json:"httpStatus"`
	Message              string `json:"message"`
	InternalErrorMessage string `json:"internalErrorMessage,omitempty"`
}

func NewAPIError(id int, httpStatus int, errMsg string, err error) *APIError {
	apiError := &APIError{
		ErrorID:    id,
		HTTPStatus: httpStatus,
		Message:    errMsg,
	}
	if err != nil {
		apiError.InternalErrorMessage = err.Error()
	}
	return apiError
}

// NewAPI404Error - Resource not found
func NewAPI404Error(id int, errMsg string, err error) *APIError {
	apiError := &APIError{
		ErrorID:    id,
		HTTPStatus: http.StatusNotFound,
		Message:    errMsg,
	}
	if err != nil {
		apiError.InternalErrorMessage = err.Error()
	}
	return apiError
}

// NewAPI400Error - Bad request error
func NewAPI400Error(id int, errMsg string, err error) *APIError {
	apiError := &APIError{
		ErrorID:    id,
		HTTPStatus: http.StatusBadRequest,
		Message:    errMsg,
	}
	if err != nil {
		apiError.InternalErrorMessage = err.Error()
	}
	return apiError
}

// NewAPI401Error - Unauthorized error or Not authenticated
func NewAPI401Error(id int, errMsg string, err error) *APIError {
	apiError := &APIError{
		ErrorID:    id,
		HTTPStatus: http.StatusUnauthorized,
		Message:    errMsg,
	}
	if err != nil {
		apiError.InternalErrorMessage = err.Error()
	}
	return apiError
}

// NewAPI403Error - Forbidden, client identity is known by server, but he has no authorize to access the requested resources
func NewAPI403Error(id int, errMsg string, err error) *APIError {
	apiError := &APIError{
		ErrorID:    id,
		HTTPStatus: http.StatusForbidden,
		Message:    errMsg,
	}
	if err != nil {
		apiError.InternalErrorMessage = err.Error()
	}
	return apiError
}

// NewAPI500Error - The server has encountered a situation it doesn't know how to handle.
func NewAPI500Error(id int, errMsg string, err error) *APIError {
	apiError := &APIError{
		ErrorID:    id,
		HTTPStatus: http.StatusInternalServerError,
		Message:    errMsg,
	}
	if err != nil {
		apiError.InternalErrorMessage = err.Error()
	}
	return apiError
}
