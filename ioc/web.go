package ioc

import (
	"github.com/gin-gonic/gin"
	"mall/internal/auth"
	"mall/internal/auth/jwt"
	"mall/internal/product"
	"mall/internal/user"
)

func InitWeb(mdl []gin.HandlerFunc, userHdl *user.Handler, productHdl *product.Handler) *gin.Engine {
	server := gin.Default()
	server.Use(mdl...)

	userHdl.RegisterRoute(server)
	productHdl.RegisterRoute(server)

	return server
}

func InitGinMiddlewares(jwtHdl *jwt.TokenHandler, sessionHdl *jwt.RedisSession) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		auth.CORS(),
		auth.NewAuthorityMiddlewareBuilder(jwtHdl, sessionHdl).
			IgnorePath("/api/user/signup").
			IgnorePath("/api/user/login").
			IgnorePath("/api/user/send-code").
			IgnorePath("/api/user/verify-code").
			MerchantPath("/api/products/:id/onlist").
			MerchantPath("/api/products/:id/removelist").
			MerchantPath("/api/products/:id").
			MerchantPath("/api/products/").
			MerchantPath("/api/categories/").
			CheckLogin(),
	}
}
