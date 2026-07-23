package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck/checks"
)

func (a *App) initHealth() {
	health := healthcheck.New(a.cfg.ServiceName, a.cfg.Version, a.logger)
	health.Add(
		checks.NewGRPC(a.manager),
		checks.NewRedis(a.redis),
		checks.NewSystem(),
	)

	health.AddReadiness(checks.NewGRPC(a.manager), checks.NewRedis(a.redis))
	health.AddLiveness(checks.NewSystem())

	a.health = health

}
