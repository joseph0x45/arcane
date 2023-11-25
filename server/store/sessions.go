package store

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Sessions struct {
	redis *redis.Client
}

func NewSessionsStore(redis *redis.Client) *Sessions {
	return &Sessions{
		redis: redis,
	}
}

func (s *Sessions) Create(sessionValue string) (sessionId string, err error) {
	ctx := context.Background()
	sessionId = uuid.NewString()
	err = s.redis.Set(
		ctx,
		sessionId,
		sessionValue,
		0,
	).Err()
	return
}

func (s *Sessions) GetById(id string) (sessionValue string, err error) {
	ctx := context.Background()
	sessionValue, err = s.redis.Get(
		ctx,
		id,
	).Result()
	return
}

func (s *Sessions) Delete(id string) error {
	ctx := context.Background()
	err := s.redis.Del(ctx, id).Err()
	return err
}
