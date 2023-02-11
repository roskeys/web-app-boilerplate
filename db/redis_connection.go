package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

var Rdb *redis.Client
var RedisLimiter *redis_rate.Limiter

func InitRedis() bool {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	Rdb = rdb
	limiter := redis_rate.NewLimiter(rdb)
	RedisLimiter = limiter
	return true
}

var _ = InitRedis()
