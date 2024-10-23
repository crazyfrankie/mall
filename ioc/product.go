package ioc

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"mall/internal/product"
	"mall/internal/product/web"
)

func InitProduct(db *gorm.DB, cmd redis.Cmdable) *web.ProductHandler {
	return product.InitProductHandler(db, cmd)
}
