package dto

type (
	PaginationRequest struct {
		Page     int    `form:"page,default=1" binding:"min=1"`
		PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
		Search   string `form:"search,omitempty"`
	}

	PaginationResponse[T any] struct {
		Datas      []*T  `json:"data"`
		Total      int64 `json:"total"`
		Page       int   `json:"page"`
		PageSize   int   `json:"page_size"`
		TotalPages int   `json:"total_pages"`
	}
)
