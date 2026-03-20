package auth

import (
	"net/http"
	"strings"

	"github.com/aK1r4z/workpal/internal/pkg/response"
	"github.com/aK1r4z/workpal/internal/session"
	"github.com/labstack/echo/v5"
)

const (
	ContextUserIDKey = "user_id"
)

func Middleware(sessions session.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// 获取认证 header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Fail(c, http.StatusUnauthorized, -1, "missing authorization header")
			}

			// 从 header 中提取 token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Fail(c, http.StatusBadRequest, -1, "invalid authorization format")
			}
			token := parts[1]

			// 校验 token 并获取用户标识符
			userID, err := sessions.Get(c.Request().Context(), token)
			if err != nil {
				return response.Fail(c, http.StatusUnauthorized, -1, "invalid or expired token")
			}

			// 把 userID 注入到 echo.Context 里
			c.Set(ContextUserIDKey, userID)

			return next(c)
		}
	}
}
