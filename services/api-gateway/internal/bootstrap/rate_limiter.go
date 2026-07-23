package bootstrap

import (
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/ratelimiter"
)

func (a *App) initRateLimiter() {

	a.rateLimiter = ratelimiter.NewRedisLimiter(
		a.redis,
	)
}
