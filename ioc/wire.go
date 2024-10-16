//go:build wireinject

package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var BaseSet = wire.NewSet(InitRedis, InitDB)

func InitGin() *gin.Engine {
	wire.Build(
		BaseSet,

		UserSet,

		InitGinMiddlewares,

		InitWeb,
	)
	return nil
}
