//go:build wireinject

package user

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mall/internal/auth"
	"mall/internal/user/repository"
	"mall/internal/user/repository/cache"
	"mall/internal/user/repository/dao"
	"mall/internal/user/service"
	"mall/internal/user/web"
)

var userSet = wire.NewSet(
	dao.NewUserDao,

	cache.NewUserCache,
	cache.NewCodeCache,

	repository.NewUserRepository,
	repository.NewCodeRepository,

	InitSMSService,
	service.NewUserService,
	service.NewCodeService,

	InitLogger,
	web.NewUserHandler,
)

func InitUserHandler(db *gorm.DB, cmd redis.Cmdable) *web.UserHandler {
	wire.Build(
		auth.JWTSet,

		userSet,
	)
	return new(web.UserHandler)
}
