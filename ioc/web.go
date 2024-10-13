package ioc

import (
	"github.com/gin-gonic/gin"

	"mall/web"
)

func InitWeb(userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()

	userHdl.RegisterRoute(server)

	return server
}
