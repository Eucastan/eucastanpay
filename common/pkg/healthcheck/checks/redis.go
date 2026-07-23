package checks

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{
		client: client,
	}
}

func (r *Redis) Name() string {
	return "redis"
}

func (r *Redis) Check(ctx context.Context) healthcheck.Component {

	started := time.Now()

	ctx, cancel := context.WithTimeout(
		ctx,
		2*time.Second,
	)
	defer cancel()

	if err := r.client.Ping(ctx).Err(); err != nil {
		return healthcheck.Component{
			Name:     r.Name(),
			Status:   healthcheck.Unhealthy,
			Error:    err.Error(),
			Duration: time.Since(started).String(),
		}
	}

	size := r.client.DBSize(ctx)
	stats := r.client.PoolStats()

	return healthcheck.Component{
		Name:     r.Name(),
		Status:   healthcheck.Healthy,
		Duration: time.Since(started).String(),
		Details: map[string]interface{}{
			"db_size":     size.Val(),
			"hits":        stats.Hits,
			"misses":      stats.Misses,
			"timeouts":    stats.Timeouts,
			"total_conns": stats.TotalConns,
			"idle_conns":  stats.IdleConns,
			"stale_conns": stats.StaleConns,
		},
	}
}
