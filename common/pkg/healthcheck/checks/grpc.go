package checks

import (
	"context"
	"time"

	commongrpc "github.com/Eucastan/eucastanpay/common/pkg/grpc"
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"

	"google.golang.org/grpc/connectivity"
)

type GRPC struct {
	manager *commongrpc.Manager
}

func NewGRPC(manager *commongrpc.Manager) *GRPC {
	return &GRPC{
		manager: manager,
	}
}

func (g *GRPC) Name() string {
	return "grpc"
}

func (g *GRPC) Check(ctx context.Context) healthcheck.Component {

	started := time.Now()

	details := map[string]interface{}{}
	status := healthcheck.Healthy

	for service, conn := range g.manager.All() {
		state := conn.GetState()
		details[service] = state.String()

		switch state {
		case connectivity.Ready:

		case connectivity.Idle:

		case connectivity.Connecting:
			if status == healthcheck.Healthy {
				status = healthcheck.Degraded
			}

		default:
			status = healthcheck.Unhealthy
		}
	}

	return healthcheck.Component{
		Name:     g.Name(),
		Status:   status,
		Duration: time.Since(started).String(),
		Details:  details,
	}
}
