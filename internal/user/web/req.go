package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WrapReq[T any](fn func(c *gin.Context, req T) (Response, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.Bind(&req); err != nil {
			// 打日志
			return
		}
		res, err := fn(c, req)
		if err != nil {
			// 打日志
		}

		c.JSON(http.StatusOK, res)
	}
}
