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

func (m *Manager) Add(name string, conn *grpc.ClientConn) {

	m.mu.Lock()

	defer m.mu.Unlock()

	m.connections[name] = conn
}

func (m *Manager) Get(name string) *grpc.ClientConn {

	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.connections[name]
}

func (m *Manager) Close() error {

	m.mu.Lock()

	defer m.mu.Unlock()

	for _, conn := range m.connections {
		conn.Close()
	}

	return nil
}
