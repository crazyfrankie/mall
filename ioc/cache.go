package ioc

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr         string `yaml:"addr"`
		PoolSize     int    `yaml:"poolSize"`
		MinIdleConns int    `yaml:"minIdleConns"`
	}

	var cfg Config
	err := viper.UnmarshalKey("redis", &cfg)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     "",
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  time.Minute * 5,
	})
}
