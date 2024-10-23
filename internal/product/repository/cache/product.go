package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"mall/internal/product/domain"
)

var (
	ErrProductNotFound = errors.New("product not found")
)

type ProductCache struct {
	cmd redis.Cmdable
}

func NewProductCache(cmd redis.Cmdable) *ProductCache {
	return &ProductCache{
		cmd: cmd,
	}
}

func (cache *ProductCache) SetProduct(ctx context.Context, id uint64, productDetail domain.ProductDetail) error {
	key := cache.key(id)

	val, err := json.Marshal(productDetail)
	if err != nil {
		return err
	}

	return cache.cmd.Set(ctx, key, string(val), time.Minute*5).Err()
}

func (cache *ProductCache) GetProduct(ctx context.Context, id uint64) (domain.ProductDetail, error) {
	key := cache.key(id)

	val, err := cache.cmd.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) { // redis返回的未命中错误
			return domain.ProductDetail{}, ErrProductNotFound // 定义此错误
		}
		return domain.ProductDetail{}, err
	}

	var detail domain.ProductDetail
	err = json.Unmarshal([]byte(val), &detail)

	return detail, err
}

func (cache *ProductCache) key(id uint64) string {
	return fmt.Sprintf("product:detail:%d", id)
}
