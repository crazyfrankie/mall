package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func WrapReq[T any](fn func(c *gin.Context, req T) (Response, error), errorHandler func(c *gin.Context, err error) (Response, bool)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.Bind(&req); err != nil {
			zap.L().Error("绑定信息错误", zap.Error(err))
			c.JSON(http.StatusBadRequest, GetResponse(WithStatus(http.StatusBadRequest), WithMsg("Bind error")))
			return
		}

		res, err := fn(c, req)
		if err != nil {
			// 调用特定错误处理函数
			if response, handled := errorHandler(c, err); handled {
				c.JSON(response.Status, response)
				return
			}

			// 默认错误处理
			c.JSON(http.StatusInternalServerError, GetResponse(WithStatus(http.StatusInternalServerError), WithMsg("Internal error")))
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
