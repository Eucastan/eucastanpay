package grpc

import (
	"time"

	"github.com/sony/gobreaker"
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
			MaxRequests: 5,
			Interval:    30 * time.Second,
			Timeout:     20 * time.Second,
		},
	)

	m.breakers[service] = cb

	return cb

}
