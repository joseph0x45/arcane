package store

import (
	"os"
	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client{
  client := redis.NewClient(
    &redis.Options{
      Addr: os.Getenv("REDIS_URL"),
      Password: "",
      DB: 0,
    },
  )
  return client
}
