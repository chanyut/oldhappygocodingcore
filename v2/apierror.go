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
