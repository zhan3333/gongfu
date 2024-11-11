package client

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gongfu/internal/config"
)

func NewRedis(conf *config.Config) (redis.Cmdable, error) {
	client := redis.NewClient(&redis.Options{
		Addr: conf.Redis.Addr,
		DB:   conf.Redis.DB,
	})
	return client, client.Ping(context.Background()).Err()
}
