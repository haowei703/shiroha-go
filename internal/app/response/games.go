package response

import (
	"shiroha.com/internal/app/model"
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
