package tag

type createRequest struct {
	Name string `json:"name"`
}

type getRequest struct {
	Name string `param:"name"`
}

type listRequest struct {
	Page  int32 `query:"page"`  // 页数
	Limit int32 `query:"limit"` // 每页最多项目数
}

type deleteRequest struct {
	Name string `param:"name"`
}
