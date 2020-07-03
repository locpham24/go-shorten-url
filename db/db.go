package db

import (
	"github.com/go-redis/redis"
)

type RedisDb struct {
	Client *redis.Client
}

func (r *RedisDb) Connect() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	r.Client = redisClient
}