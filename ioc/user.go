package ioc

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"mall/internal/user"
	"mall/internal/user/web"
)

func InitUser(db *gorm.DB, cmd redis.Cmdable) *web.UserHandler {
	return user.InitUserHandler(db, cmd)
}
