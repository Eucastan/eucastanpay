package grpc

import (
	"sync"

	"github.com/Eucastan/eucastanpay/common/pkg/grpc/discovery"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	gogrpc "google.golang.org/grpc"
)

type Manager struct {
	mu          sync.RWMutex
	registry    discovery.Registry
	connections map[string]*gogrpc.ClientConn
	breakers    map[string]*gobreaker.CircuitBreaker
}

func NewManager(registry discovery.Registry) *Manager {

	return &Manager{
		registry:    registry,
		connections: make(map[string]*gogrpc.ClientConn),
		breakers:    make(map[string]*gobreaker.CircuitBreaker),
	}
}

func (m *Manager) All() map[string]*grpc.ClientConn {

	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*grpc.ClientConn)

	for name, conn := range m.connections {
		result[name] = conn
	}

	return result
}

func (m *Manager) Add(name string, conn *grpc.ClientConn) {

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.connections[name]; exists {
		panic("grpc connection already registered: " + name)
	}

	m.connections[name] = conn
}

func (m *Manager) Get(name string) *grpc.ClientConn {

	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.connections[name]
}

func (m *Manager) Exists(name string) bool {

	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.connections[name]

	return ok
}

func (m *Manager) Close() error {

	m.mu.Lock()
	defer m.mu.Unlock()

	var firstErr error

	for name, conn := range m.connections {
		if conn == nil {
			continue
		}

		if err := conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}

		delete(m.connections, name)
	}

	return firstErr
}
