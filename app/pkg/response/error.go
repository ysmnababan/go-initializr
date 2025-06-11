package response

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	// "github.com/rs/zerolog/log"
)

type errorResponse struct {
	Meta  Meta      `json:"meta"`
	Error string    `json:"error"`
	Data  emptyData `json:"data"` // show default empty object if error occur
}

type emptyData struct{}

type Error struct {
	Response     errorResponse `json:"response"`
	Code         int           `json:"code"`
	ErrorMessage error
	Additional   any `json:"additional,omitempty"`
}

var (
	// BadRequest
	ErrBadRequest                       = CustomError(http.StatusBadRequest, 40001, "Bad Request")
	ErrValidation                       = CustomError(http.StatusBadRequest, 40002, "Invalid parameters or payload")
	ErrValueMismatch                    = CustomError(http.StatusBadRequest, 40003, "Requested value mismatch")
	ErrForgotPasswordResendTokenInvalid = CustomError(http.StatusBadRequest, 40005, "Invalid forgot password resend token")
	ErrInvalidUpdateStatus              = CustomError(http.StatusBadRequest, 40006, "Invalid status update")

	// Unauthorized
	ErrUnauthorized           = CustomError(http.StatusUnauthorized, 40101, "Unauthorized, please login")
	ErrInvalidUserCredentials = CustomError(http.StatusUnauthorized, 40102, "Invalid user credentials")
	ErrInvalidApplicationID   = CustomError(http.StatusUnauthorized, 40103, "Invalid application id")
	ErrExpiredAccessToken     = CustomError(http.StatusUnauthorized, 40104, "Access token has expired")
	ErrInvalidUserAccount     = CustomError(http.StatusUnauthorized, 40105, "Invalid user account")
	ErrInvalidPublicAuth      = CustomError(http.StatusUnauthorized, 40106, "Invalid public auth")

	// Forbidden
	ErrForbidden              = CustomError(http.StatusForbidden, 40301, "Forbidden")
	ErrAccountNotInWhitelist  = CustomError(http.StatusForbidden, 40302, "User not in whitelist")
	ErrForbiddenRoom          = CustomError(http.StatusForbidden, 40303, "You are not authorized to access this room.")
	ErrForbiddenApiPermission = CustomError(http.StatusForbidden, 40304, "Anda tidak memiliki akses resource ini")

	// NotFound
	ErrNotFound             = CustomError(http.StatusNotFound, 40401, "Data not found")
	ErrRouteNotFound        = CustomError(http.StatusNotFound, 40402, "Route not found")
	ErrAccountNotRegistered = CustomError(http.StatusNotFound, 40403, "This account is not registered. Kindly contact your administrator for assistance")

	// Conflict
	ErrDuplicate         = CustomError(http.StatusConflict, 40901, "Created value already exists")
	ErrAccountRegistered = CustomError(http.StatusConflict, 40902, "Your account has already been registered")
	ErrOrderDetailChange = CustomError(http.StatusConflict, 40903, "Order detail has been changed")

	// UnprocessableEntity
	ErrUnprocessableEntity = CustomError(http.StatusUnprocessableEntity, 42201, "Invalid parameters or payload")
	ErrInvalidPrice        = CustomError(http.StatusUnprocessableEntity, 42202, "Oops! The value entered is not a multiple of 1000. Please correct it.")
	ErrInvalidOTP          = CustomError(http.StatusUnprocessableEntity, 42205, "Invalid OTP")
	ErrInvalidToken        = CustomError(http.StatusUnprocessableEntity, 42206, "Invalid Token")

	// Too Many Request
	ErrForgotPasswordMaxAttempt = CustomError(http.StatusTooManyRequests, 40004, "You have reached the maximum request limit for today")

	// InternalServerError
	ErrInternalServerError = CustomError(http.StatusInternalServerError, 50001, "Something bad happened")
	ErrOpeningFile         = CustomError(http.StatusInternalServerError, 50002, "Error opening folder")
	ErrScanner             = CustomError(http.StatusInternalServerError, 50003, "Error bufio scanner")
	ErrCreatingFolder      = CustomError(http.StatusInternalServerError, 50004, "Error creating folder")
	ErrCreatingFile        = CustomError(http.StatusInternalServerError, 50005, "Error creating file")
	ErrExecutingTemplate   = CustomError(http.StatusInternalServerError, 50006, "Error executing template")
)

func ErrorBuilder(res *Error, message error) *Error {
	res.ErrorMessage = message
	return res
}

// ErrorWrap wrap err whit new Error base on provided errBase.
func ErrorWrap(errBase *Error, err error, detail ...any) *Error {
	if errBase == nil {
		errBase = ErrInternalServerError
	}

	resp := &Error{
		Response: errorResponse{
			Meta: Meta{
				Success:    false,
				Message:    errBase.Response.Meta.Message,
				StatusCode: errBase.Response.Meta.StatusCode,
			},
			Error: errBase.Response.Error,
		},
		Code:         errBase.Code,
		ErrorMessage: err,
	}

	if len(detail) != 0 {
		resp.Response.Meta.Detail = detail[0]
	}

	if errBase.Code != http.StatusInternalServerError {
		resp.Response.Meta.Message = err.Error()
	}

	return resp
}

func CustomError(httpCode, statusCode int, message string) *Error {
	return &Error{
		Response: errorResponse{
			Meta: Meta{
				Success:    false,
				Message:    message,
				StatusCode: statusCode,
			},
		},
		Code: httpCode,
	}
}

// CustomErrorMessage return new error
func CustomErrorMessage(base *Error, message string, err error, additional ...any) *Error {
	return &Error{
		Response: errorResponse{
			Meta: Meta{
				Success:    false,
				Message:    message,
				StatusCode: base.Response.Meta.StatusCode,
			},
			Error: base.Response.Error,
		},
		Code:         base.Code,
		ErrorMessage: err,
		Additional:   additional,
	}
}

// ErrorMessageFrom is shorthand from CustomErrorMessage whit message from
// err.Error()
func ErrorMessageFrom(base *Error, err error) *Error {
	return CustomErrorMessage(base, err.Error(), err)
}

func ErrorResponse(err error) *Error {
	re, ok := err.(*Error)
	if ok {
		return re
	} else {
		return ErrorWrap(ErrInternalServerError, err)
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("error code %d", e.Code)
}

func (e *Error) ParseToError() error {
	return e
}

func (e *Error) Send(c echo.Context, request ...any) error {
	// event := log.Error().
	// 	Stack().Err(e.ErrorMessage).
	// 	Str("uri", c.Request().RequestURI)

	// if len(request) != 0 {
	// 	event.Interface("request", request[0])
	// }

	// event.Send()

	return c.JSON(e.Code, e.Response)
}
