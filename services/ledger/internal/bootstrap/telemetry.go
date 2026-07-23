package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"go.opentelemetry.io/otel"
)

func (a *App) initTelemetry() error {
	tracer := otel.Tracer(a.cfg.ServiceName)
	meter := otel.Meter(a.cfg.ServiceName)

	tm, err := telemetry.New(tracer, meter, a.logger)
	if err != nil {
		return err
	}

	a.telemetry = tm
	return nil
}
