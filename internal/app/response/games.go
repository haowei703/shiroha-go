package response

import (
	"github.com/haowei703/shiroha/internal/app/model"
)

type Pagination struct {
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	Pages       int   `json:"pages"`
	TotalCount  int64 `json:"totalCount"`
}

type PaginatedQueryResponse struct {
	Pagination Pagination   `json:"pagination"`
	Games      []model.Game `json:"games"`
}
