package core

import (
	"errors"

	"github.com/revel/revel"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

// RenderGenericRevelJSONSuccess is wrapper function for rendering json in type of JSONResponse
func RenderGenericRevelJSONSuccess(c *revel.Controller, data interface{}) revel.Result {
	return c.RenderJSON(JSONResponse{
		Success: true,
		Data:    data,
		Errors:  nil,
	})
}

// RenderGenericRevelJSONError is wrapper function for rendering json in type of JSONResponse
func RenderGenericRevelJSONError(c *revel.Controller, err error) revel.Result {
	return c.RenderJSON(JSONResponse{
		Success: false,
		Data:    nil,
		Errors:  []string{err.Error()},
	})
}

// RenderGenericRevelJSONErrors is wrapper function for rendering json in type of JSONResponse
func RenderGenericRevelJSONErrors(c *revel.Controller, errors []error) revel.Result {
	errMessages := make([]string, 0)
	for _, err := range errors {
		errMessages = append(errMessages, err.Error())
	}
	return c.RenderJSON(JSONResponse{
		Success: false,
		Data:    nil,
		Errors:  errMessages,
	})
}

// GenericErrorFromValidationErrors is a convenient function wrapping revel.ValidationError into []error
func GenericErrorFromValidationErrors(validationErrors []*revel.ValidationError) []error {
	errorList := make([]error, 0)
	for _, e := range validationErrors {
		errorList = append(errorList, errors.New(e.Message))
	}
	return errorList
}
