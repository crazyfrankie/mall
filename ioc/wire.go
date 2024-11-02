//go:build wireinject

package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"mall/internal/auth"
	"mall/internal/product"
	"mall/internal/user"
)

var BaseSet = wire.NewSet(InitRedis, InitDB)

func InitGin() *gin.Engine {
	wire.Build(
		BaseSet,

		InitLogger,

		auth.JWTSet,

		user.InitUserHandler,

		product.InitProductHandler,

		InitMiddleware,

		InitWeb,
	)
	return nil
}
