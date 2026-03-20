package tag

import (
	"errors"
	"net/http"

	"github.com/aK1r4z/workpal/internal/auth"
	"github.com/aK1r4z/workpal/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type handler struct {
	tagService service
}

func NewHandler(tagService service) *handler {
	return &handler{
		tagService: tagService,
	}
}

// 注册路由
func (h *handler) RegisterRoutes(api *echo.Group) {
	api.POST("/tags", h.create)

	api.GET("/tags", h.list)
	api.GET("/tags/:name", h.get)

	api.DELETE("/tags/:name", h.delete)
}

// 标签创建请求处理器
func (h *handler) create(c *echo.Context) error {
	userID, ok := c.Get(auth.CtxUserID).(uuid.UUID)
	if !ok {
		// how could this even happen?
		return response.ErrUnauthorized(c)
	}

	var req createRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	err := h.tagService.Create(c.Request().Context(), userID, req.Name)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, -1, err.Error())
	}

	return nil
}

// 获取标签列表请求处理器
func (h *handler) list(c *echo.Context) error {
	userID, ok := c.Get(auth.CtxUserID).(uuid.UUID)
	if !ok {
		// you know what im going to say.
		return response.ErrUnauthorized(c)
	}

	var req listRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	// 修正参数
	if req.Page < 0 {
		req.Page = 0
	}

	if req.Limit < 0 {
		req.Limit = 5
	}

	// 转换成 ListFilter
	filter := ListFilter{
		Page:  req.Page,
		Limit: req.Limit,
	}

	tags, err := h.tagService.List(c.Request().Context(), userID, filter)
	if err != nil {
		return response.Fail(c, http.StatusInternalServerError, -1, err.Error())
	}

	resp := listResponse{
		Tags: make([]string, len(tags)),
	}

	for i, t := range tags {
		resp.Tags[i] = t.Name
	}

	return response.Success(c, http.StatusOK, 0, "success", resp)
}

// 标签信息获取请求处理器
func (h *handler) get(c *echo.Context) error {
	userID, ok := c.Get(auth.CtxUserID).(uuid.UUID)
	if !ok {
		// Still, i think this will never happen, but who knows?
		return response.ErrUnauthorized(c)
	}

	var req getRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	t, err := h.tagService.Get(c.Request().Context(), userID, req.Name)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return response.Success(c, http.StatusOK, 0, "not found", response.Nil{})
		}
		return response.ErrInternalServerError(c, err)
	}

	return response.Success(c, http.StatusOK, 0, "found", getResponse{
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
	})
}

// 删除标签请求处理器
func (h *handler) delete(c *echo.Context) error {
	userID, ok := c.Get(auth.CtxUserID).(uuid.UUID)
	if !ok {
		// well... yeah.
		return response.ErrUnauthorized(c)
	}

	var req deleteRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	err := h.tagService.Delete(c.Request().Context(), userID, req.Name)
	if err != nil {
		return response.ErrInternalServerError(c, err)
	}

	return response.Success(c, http.StatusOK, 0, "deleted", response.Nil{})
}
