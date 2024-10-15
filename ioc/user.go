package ioc

import (
	"github.com/google/wire"
	"mall/middleware/jwt"
	"mall/repository"
	"mall/repository/cache"
	"mall/repository/dao"
	"mall/service"
	"mall/web"
)

var JWTSet = wire.NewSet(
	jwt.NewJwtHandler,
	jwt.NewRedisSession,
)

var UserSet = wire.NewSet(
	dao.NewUserDao,

	cache.NewUserCache,
	cache.NewCodeCache,

	repository.NewUserRepository,
	repository.NewCodeRepository,

	InitSMSService,
	service.NewUserService,
	service.NewCodeService,

	JWTSet,

	web.NewUserHandler,
)
