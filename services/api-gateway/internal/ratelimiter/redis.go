package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(client *redis.Client) *RedisLimiter {
	return &RedisLimiter{
		client: client,
	}
}

func (r *RedisLimiter) Allow(ctx context.Context, key string, limit int, windowSeconds int) (bool, error) {

	current, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if current == 1 {

		err = r.client.Expire(
			ctx,
			key,
			time.Duration(windowSeconds)*time.Second,
		).Err()

		if err != nil {
			return false, err
		}
	}

	if current > int64(limit) {
		return false, nil
	}

	return true, nil
}
