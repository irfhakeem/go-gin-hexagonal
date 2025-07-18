package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

func Success(c *gin.Context, message string, data any, code int) {
	c.JSON(code, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, message string, data any, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Status:  true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *gin.Context, message string, err string, code int) {
	c.JSON(code, Response{
		Status:  false,
		Message: message,
		Error:   err,
	})
}
