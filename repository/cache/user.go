package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mall/domain"
	"time"
)

type UserCache struct {
	cmd redis.Cmdable
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd: cmd,
	}
}

func (cache *UserCache) Set(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}

	key := cache.key(user.Id)

	return cache.cmd.Set(ctx, key, val, time.Minute*10).Err()
}

func (cache *UserCache) Get(ctx context.Context, id uint64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}

	var user domain.User
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}

func (cache *UserCache) key(id uint64) string {
	return fmt.Sprintf("user:%d:info", id)
}
