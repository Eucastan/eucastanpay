package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck/checks"
)

func (a *App) initHealth() {
	h := healthcheck.New(a.cfg.ServiceName, a.cfg.Version, a.logger)
	h.Add(
		checks.NewDatabase(a.database.DB),
		checks.NewGRPC(a.manager),
		checks.NewKafkaProducer(a.publish),
		checks.NewSystem(),
	)

	h.AddReadiness(
		checks.NewDatabase(a.database.DB),
		checks.NewGRPC(a.manager),
		checks.NewKafkaProducer(a.publish),
	)

	h.AddLiveness(checks.NewSystem())

	a.health = h
}
