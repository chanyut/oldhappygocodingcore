package core

import (
	"fmt"
	"net/http"
	"strings"

	"mime/multipart"

	"io/ioutil"

	"github.com/revel/revel"
)

type RevelRequestMethodType string

const (
	RevelRequestMethodTypeGET  RevelRequestMethodType = "GET"
	RevelRequestMethodTypePOST                        = "POST"
	RevelRequestMethodTypePUT                         = "PUT"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   *APIError   `json:"error"`
}

func (resp JSONResponse) InternalErrorString() string {
	if resp.Error == nil {
		return ""
	}
	return resp.Error.InternalErrorMessage
}

type RevelResultRenderer struct {
	controller *revel.Controller
}

func NewRevelResultRenderer(controller *revel.Controller) RevelResultRenderer {
	if controller == nil {
		panic(fmt.Errorf("controller cannot be null"))
	}
	renderer := RevelResultRenderer{
		controller: controller,
	}
	return renderer
}

// RenderJSONSuccess is wrapp er function for rendering json in type of JSONResponse
func (r *RevelResultRenderer) RenderJSONSuccess(data interface{}) revel.Result {
	return r.controller.RenderJSON(JSONResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	})
}

// RenderJSONError is wrapper function for rendering json in type of JSONResponse
func (r *RevelResultRenderer) RenderJSONError(err *APIError) revel.Result {
	if err == nil {
		panic(fmt.Errorf("err must not be null"))
	}

	if err.HTTPStatus == 0 {
		err.HTTPStatus = http.StatusInternalServerError
	}
	r.controller.Response.Status = err.HTTPStatus
	return r.controller.RenderJSON(JSONResponse{
		Success: false,
		Data:    nil,
		Error:   err,
	})
}

// RenderJSONError is wrapper function for rendering json in type of JSONResponse
func (r *RevelResultRenderer) RenderJSONErrorWithUnspecificError(err error) revel.Result {
	if err == nil {
		panic(fmt.Errorf("err must not be null"))
	}

	apiError := NewAPI500Error(0, err.Error(), nil)
	r.controller.Response.Status = apiError.HTTPStatus
	return r.controller.RenderJSON(JSONResponse{
		Success: false,
		Data:    nil,
		Error:   apiError,
	})
}

// RenderCurrentValidationError is wrapper function for rendering json in type of JSONResponse
func (r *RevelResultRenderer) RenderCurrentValidationError() revel.Result {
	errMessages := ""
	for _, err := range r.controller.Validation.Errors {
		errMessages += fmt.Sprintf("%s: %s\n", err.Key, err.Message)
	}

	validationError := NewAPI400Error(0, errMessages, nil)
	r.controller.Response.Status = http.StatusBadRequest
	return r.controller.RenderJSON(JSONResponse{
		Success: false,
		Data:    nil,
		Error:   validationError,
	})
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

func GetMethodTypeOfRevelRequest(req *revel.Request) RevelRequestMethodType {
	if strings.ToLower(req.Method) == "get" {
		return RevelRequestMethodTypeGET
	} else if strings.ToLower(req.Method) == "post" {
		return RevelRequestMethodTypePOST
	} else if strings.ToLower(req.Method) == "put" {
		return RevelRequestMethodTypePUT
	}
	panic(fmt.Errorf("unknow request's method: %v", req.Method))
}
