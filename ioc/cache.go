package ioc

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		MinIdleConns: 5,
		PoolSize:     15,
		DialTimeout:  time.Minute * 5,
	})
}
