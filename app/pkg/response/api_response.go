package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Meta struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"` //internal error code
	Detail     any    `json:"detail,omitempty"`
}

type APIResponse struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data"`
}

// func WithErrorResponse(err error, c echo.Context) error {
// 	var meta Meta
// 	re, ok := err.(*APIError)
// 	meta.Success = false
// 	code := re.Code // http code
// 	if ok {
// 		meta.Message = re.Message
// 		meta.StatusCode = re.StatusCode
// 		meta.Detail = re.Detail
// 	} else {
// 		defErr := ErrInternalServerError
// 		meta.Message = defErr.Message
// 		meta.StatusCode = defErr.StatusCode
// 		code = http.StatusInternalServerError
// 	}

// 	return c.JSON(code, APIResponse{Meta: meta})
// }

func WithStatusOKResponse(data any, c echo.Context) error {
	resp := APIResponse{
		Meta: Meta{
			Success:    true,
			Message:    "Request processed successfully",
			StatusCode: 20001,
		},
		Data: data,
	}
	return c.JSON(http.StatusOK, resp)
}

func WithStatusCreatedResponse(data any, c echo.Context) error {
	resp := APIResponse{
		Meta: Meta{
			Success:    true,
			Message:    "Request created successfully",
			StatusCode: 20101,
		},
		Data: data,
	}
	return c.JSON(http.StatusCreated, resp)
}
