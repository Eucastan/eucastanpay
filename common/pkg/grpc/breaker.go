package grpc

import (
	"context"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/errmessage"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

func (m *Manager) breaker(service string) *gobreaker.CircuitBreaker {

	m.mu.Lock()

	defer m.mu.Unlock()

	if cb, ok := m.breakers[service]; ok {
		return cb
	}

	cb := gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        service,
			MaxRequests: 3,
			Interval:    60 * time.Second,
			Timeout:     30 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {

				if counts.Requests < 10 {
					return false
				}

				failureRate := float64(counts.TotalFailures) /
					float64(counts.Requests)

				return failureRate >= 0.6
			},
		},
	)

	m.breakers[service] = cb

	return cb

}

func (m *Manager) Invoke(
	ctx context.Context,
	service string,
	method string,
	req any,
	resp any,
	opts ...grpc.CallOption,
) error {

	conn := m.Get(service)
	if conn == nil {
		return errmessage.ErrServiceNotFound
	}

	cb := m.breaker(service)

	_, err := cb.Execute(func() (interface{}, error) {
		err := conn.Invoke(ctx, method, req, resp, opts...)
		return nil, err
	})

	return err
}
