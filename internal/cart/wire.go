//go:build wireinject

package cart

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mall/internal/cart/repository"
	"mall/internal/cart/repository/dao"
	"mall/internal/cart/service"
	"mall/internal/cart/web"
	"mall/internal/product"
)

func NewCartHandler(db *gorm.DB, cmd redis.Cmdable) *web.CartHandler {
	wire.Build(
		dao.NewCartDao,

		repository.NewCartRepository,

		product.NewProductRepository,

		service.NewCartService,

		web.NewCartHandler,
	)
	return new(web.CartHandler)
}
