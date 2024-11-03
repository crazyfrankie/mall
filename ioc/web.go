package ioc

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"mall/internal/auth"
	"mall/internal/auth/jwt"
	"mall/internal/product"
	"mall/internal/user"
	logger2 "mall/pkg/logger"
	"mall/pkg/middleware/logger"
)

func InitWeb(mdl []app.HandlerFunc, userHdl *user.Handler, productHdl *product.Handler) *server.Hertz {
	server := server.Default()
	server.Use(mdl...)
	userHdl.RegisterRoute(server)
	productHdl.RegisterRoute(server)

	return server
}

func InitMiddleware(jwtHdl *jwt.TokenHandler, sessionHdl *jwt.RedisSession, l logger2.Logger) []app.HandlerFunc {
	return []app.HandlerFunc{
		auth.NewTokenEffectiveBuilder(jwtHdl, sessionHdl).
			IgnorePath("/api/user/send-code").
			IgnorePath("/api/user/verify-code").
			IgnorePath("/api/user/login").
			Check(),

		logger.NewMiddlewareBuilder(func(ctx context.Context, al *logger.AccessLog) {
			l.Debug("HTTP 请求", logger2.Field{Key: "req", Val: al})
		}).AllowReqBody().AllowRespBody().Build(),
	}
}
