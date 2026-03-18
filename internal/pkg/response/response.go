package response

import (
	"github.com/labstack/echo/v5"
)

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
