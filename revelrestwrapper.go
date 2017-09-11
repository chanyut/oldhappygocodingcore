package core

import (
	"fmt"

	"mime/multipart"

	"io/ioutil"

	"github.com/revel/revel"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

func (resp JSONResponse) ErrorString() string {
	if resp.Errors == nil {
		return ""
	} else if len(resp.Errors) == 0 {
		return ""
	}

	errorString := ""
	for _, errStr := range resp.Errors {
		errorString += errStr + "\n"
	}
	return errorString
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
	// log.Printf("[RevelRestWrapper] return error: %v", err)
	if err == nil {
		panic(fmt.Errorf("err must not be null"))
	}
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
	// log.Printf("[RevelRestWrapper] return error: %v", errors)
	return c.RenderJSON(JSONResponse{
		Success: false,
		Data:    nil,
		Errors:  errMessages,
	})
}

// RenderGenericRevelJSONWithValidationErrors is wrapper function for rendering json in type of JSONResponse
func RenderGenericRevelJSONWithValidationErrors(c *revel.Controller, errors []*revel.ValidationError) revel.Result {
	errMessages := make([]string, 0)
	for _, err := range errors {
		errMessages = append(errMessages, err.Message)
	}
	// log.Printf("[RevelRestWrapper] return error: %v", errors)
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
		errorList = append(errorList, fmt.Errorf("%s - %s", e.Key, e.Message))
	}
	return errorList
}

// GetFileDataFromFileHeader read data from given fileHeader (multipart.FileHeader) and return as []byte
func GetFileDataFromFileHeader(fileHeader *multipart.FileHeader) ([]byte, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}
