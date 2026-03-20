package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Nil struct{}

type Response[T any] struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func Success[T any](c *echo.Context, httpStatus int, code int32, message string, data T) error {
	return c.JSON(httpStatus, Response[T]{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Fail(c *echo.Context, httpStatus int, code int32, message string) error {
	return c.JSON(httpStatus, Response[any]{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// Response Message
const (
	MsgInvalidRequest = "invalid request"
	MsgInvalidToken   = "invalid or expired token"
)

// Errors
func ErrBadRequest(c *echo.Context) error {
	return Fail(c, http.StatusBadRequest, -1, MsgInvalidRequest)
}

func ErrUnauthorized(c *echo.Context) error {
	return Fail(c, http.StatusUnauthorized, -1, MsgInvalidToken)
}

func ErrInternalServerError(c *echo.Context, err error) error {
	return Fail(c, http.StatusInternalServerError, -1, err.Error())
}
