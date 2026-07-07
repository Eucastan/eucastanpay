package telemetry

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	tracer trace.Tracer
	meter  metric.Meter
	logger *logrus.Logger

	RequestCounter metric.Int64Counter
	RequestLatency metric.Float64Histogram
}

func New(tracer trace.Tracer, meter metric.Meter, logger *logrus.Logger) (*Telemetry, error) {
	counter, err := meter.Int64Counter("http.requests.total")
	if err != nil {
		return nil, err
	}

	latency, err := meter.Float64Histogram("http.request.duration")
	if err != nil {
		return nil, err
	}

	return &Telemetry{
		tracer: tracer,
		meter:  meter,
		logger: logger,

		RequestCounter: counter,
		RequestLatency: latency,
	}, nil
}

func (t *Telemetry) Tracer() trace.Tracer {
	return t.tracer
}

func (t *Telemetry) Meter() metric.Meter {
	return t.meter
}

func (t *Telemetry) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {

	return t.tracer.Start(ctx, name, opts...)
}

func (t *Telemetry) RecordError(
	span trace.Span,
	err error,
) {
	if err == nil {
		return
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

func (t *Telemetry) Logger(ctx context.Context) *logrus.Entry {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()

	return t.logger.WithFields(logrus.Fields{
		"trace_id": sc.TraceID().String(),
		"span_id":  sc.SpanID().String(),
	})
}

func (t *Telemetry) Info(args ...interface{}) {
	t.logger.Info(args...)
}

func (t *Telemetry) Error(args ...interface{}) {
	t.logger.Error(args...)
}
