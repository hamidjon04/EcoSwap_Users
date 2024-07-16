package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserRedis struct {
	Rdb *redis.Client
}

func NewRedisRepo(rdb *redis.Client) *UserRedis {
	return &UserRedis{
		Rdb: rdb,
	}
}

func (U *UserRedis) AddToBlacklist(ctx context.Context, token string, duration time.Duration) error {
	err := U.Rdb.Set(ctx, token, "blacklisted", duration).Err()
	return err
}
