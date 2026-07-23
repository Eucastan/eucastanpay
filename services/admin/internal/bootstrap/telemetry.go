package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"go.opentelemetry.io/otel"
)

func (a *App) initTelemetry() error {
	tracer := otel.Tracer("admin-service")
	meter := otel.Meter("admin-service")

	tm, err := telemetry.New(tracer, meter, a.logger)
	if err != nil {
		return err
	}

	a.telemetry = tm

	return nil
}
