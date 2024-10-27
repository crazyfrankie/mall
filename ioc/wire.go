//go:build wireinject

package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"mall/internal/auth"
)

var BaseSet = wire.NewSet(InitRedis, InitDB)

func InitGin() *gin.Engine {
	wire.Build(
		BaseSet,

		InitLogger,

		auth.JWTSet,

		InitUser,

		InitProduct,

		InitMiddleware,

		InitWeb,
	)
	return nil
}
