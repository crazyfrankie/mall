package ioc

import (
	"github.com/gin-gonic/gin"
	"mall/internal/auth"
	"mall/internal/auth/jwt"
	"mall/internal/product"
	"mall/internal/user"
)

func InitWeb(mdl gin.HandlerFunc, userHdl *user.Handler, productHdl *product.Handler) *gin.Engine {
	server := gin.Default()
	server.Use(mdl)
	userHdl.RegisterRoute(server)
	productHdl.RegisterRoute(server)

	return server
}

func InitMiddleware(jwtHdl *jwt.TokenHandler, sessionHdl *jwt.RedisSession) gin.HandlerFunc {
	return auth.NewTokenEffectiveBuilder(jwtHdl, sessionHdl).
		IgnorePath("/api/user/send-code").
		IgnorePath("/api/user/verify-code").
		IgnorePath("/api/user/login").
		Check()
}
