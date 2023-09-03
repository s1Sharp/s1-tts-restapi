package storage

import (
	"context"
	"github.com/go-redis/redis"
	"log"
)

type RedisStorage struct {
	ctx    context.Context
	Client *redis.Client
}

func NewRedisStorage(uri string, healthcheckMessage string, ctx context.Context) (rs RedisStorage) {
	log.Println("Redis database connecting...")

	// Connect to Redis
	opt, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(opt)
	if _, err := redisClient.Ping().Result(); err != nil {
		panic(err)
	}

	err = redisClient.Set("healthcheck", healthcheckMessage, 0).Err()
	if err != nil {
		panic(err)
	}

	rs = RedisStorage{
		ctx:    ctx,
		Client: redisClient,
	}

	log.Println("Redis database Connected.")
	return
}
