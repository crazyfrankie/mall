package ioc

import (
	"github.com/google/wire"

	"mall/internal/user/middleware/jwt"
	"mall/internal/user/repository"
	"mall/internal/user/repository/cache"
	"mall/internal/user/repository/dao"
	"mall/internal/user/service"
	"mall/internal/user/web"
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
