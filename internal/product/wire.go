//go:build wireinject

package product

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mall/internal/product/repository"
	"mall/internal/product/repository/cache"
	"mall/internal/product/repository/dao"
	"mall/internal/product/service"
	"mall/internal/product/web"
)

var productSet = wire.NewSet(
	dao.NewProductDao,
	cache.NewProductCache,

	repository.NewProductRepository,

	service.NewProductService,

	web.NewProductHandler,
)

func InitProductHandler(db *gorm.DB, cmd redis.Cmdable) *web.ProductHandler {
	wire.Build(
		productSet,
	)
	return new(web.ProductHandler)
}

func NewProductRepository(db *gorm.DB, cmd redis.Cmdable) *repository.ProductRepository {
	wire.Build(
		dao.NewProductDao,
		cache.NewProductCache,
		repository.NewProductRepository,
	)
	return new(repository.ProductRepository)
}
