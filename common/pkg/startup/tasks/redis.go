package tasks

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Run(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
