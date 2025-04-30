package dao

import (
	"context"
	"ip_pool_chick/config"

	"github.com/go-redis/redis/v8"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis(cfg config.RedisConfig) {
	Rdb = redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
}
