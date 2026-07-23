package bootstrap

import "github.com/Eucastan/eucastanpay/services/user/internal/infra/redis"

func (a *App) initRedis() {
	a.redis = redis.NewRedisClient(a.cfg, a.logger)
}
