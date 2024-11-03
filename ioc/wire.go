//go:build wireinject

package ioc

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	"mall/internal/auth"
)

var BaseSet = wire.NewSet(InitRedis, InitDB)

func InitGin() *server.Hertz {
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
