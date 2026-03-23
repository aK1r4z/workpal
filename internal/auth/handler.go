package auth

import (
	"net/http"

	"github.com/aK1r4z/workpal/pkg/response"
	"github.com/labstack/echo/v5"
)

type handler struct {
	authService *service
}

// 创建认证请求处理器
func NewHandler(
	authService *service,
) *handler {
	return &handler{
		authService: authService,
	}
}

// 注册路由
func (h *handler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/auth")
	g.POST("/register", h.register)
	g.POST("/login", h.login)
}

// 注册请求处理器
// [TODO] RateLimiter
func (h *handler) register(c *echo.Context) error {
	req := &registerRequest{}
	if err := c.Bind(req); err != nil {
		return response.ErrBadRequest(c)
	}

	err := h.authService.Register(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, -1, err.Error())
	}

	return response.Success(c, http.StatusCreated, 0, "user created", true)
}

// 登录请求处理器
// [TODO] RateLimiter
func (h *handler) login(c *echo.Context) error {
	req := &loginRequest{}
	if err := c.Bind(req); err != nil {
		return response.ErrBadRequest(c)
	}

	token, err := h.authService.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return response.Fail(c, http.StatusUnauthorized, -1, "invalid credential")
	}

	resp := &loginResponse{
		Token: token,
	}

	return response.Success(c, http.StatusOK, 0, "login success", resp)
}
