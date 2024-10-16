package ioc

import (
	"github.com/gin-gonic/gin"
	"mall/middleware/jwt"
	"mall/web/auth"

	"mall/web"
)

func InitWeb(mdl []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdl...)

	userHdl.RegisterRoute(server)

	return server
}

func InitGinMiddlewares(jwtHdl *jwt.TokenHandler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		auth.CORS(),
		auth.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePath("/api/user/signup").
			IgnorePath("/api/user/login").
			IgnorePath("/api/user/send-code").
			IgnorePath("/api/user/signup/verify-code").
			IgnorePath("/api/user/login/verify-code").
			CheckLogin(),
	}
}
