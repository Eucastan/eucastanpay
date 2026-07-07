package telemetry

import (
	"context"
)

func (t *Telemetry) Count(ctx context.Context, value int64) {
	t.RequestCounter.Add(ctx, value)
}

func (t *Telemetry) ObserveLatency(ctx context.Context, ms float64) {
	t.RequestLatency.Record(ctx, ms)
}
