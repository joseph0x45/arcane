package store

import (
	"github.com/redis/go-redis/v9"
	"os"
)

func GetRedisClient() *redis.Client {
	options, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(
    options,
	)
	return client
}
