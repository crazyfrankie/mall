package ioc

import (
	"github.com/google/wire"

	"mall/repository"
	"mall/repository/dao"
	"mall/service"
	"mall/web"
)

var UserSet = wire.NewSet(
	dao.NewUserDao,

	repository.NewUserRepository,

	service.NewUserService,

	web.NewUserHandler,
)
