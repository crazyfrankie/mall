// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package user

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mall/internal/auth/jwt"
	"mall/internal/user/repository"
	"mall/internal/user/repository/cache"
	"mall/internal/user/repository/dao"
	"mall/internal/user/service"
	"mall/internal/user/web"
)

// Injectors from wire.go:

func InitUserHandler(db *gorm.DB, cmd redis.Cmdable) *web.UserHandler {
	userDao := dao.NewUserDao(db)
	userRepository := repository.NewUserRepository(userDao)
	redisSession := jwt.NewRedisSession(cmd)
	userService := service.NewUserService(userRepository, redisSession)
	codeCache := cache.NewCodeCache(cmd)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := InitSMSService()
	codeService := service.NewCodeService(codeRepository, smsService)
	tokenHandler := jwt.NewJwtHandler()
	userHandler := web.NewUserHandler(userService, codeService, tokenHandler)
	return userHandler
}

// wire.go:

var userSet = wire.NewSet(dao.NewUserDao, cache.NewUserCache, cache.NewCodeCache, repository.NewUserRepository, repository.NewCodeRepository, InitSMSService, service.NewUserService, service.NewCodeService, web.NewUserHandler)