package tasks

import (
	"context"
	"fmt"

	"github.com/Eucastan/eucastanpay/common/pkg/grpc"
	"google.golang.org/grpc/connectivity"
)

type GRPC struct {
	manager *grpc.Manager
}

func (g *GRPC) Run(ctx context.Context) error {

	for name, conn := range g.manager.All() {

		if conn == nil {
			return fmt.Errorf("%s connection is nil", name)
		}

		state := conn.GetState()

		if state != connectivity.Ready &&
			state != connectivity.Idle {

			return fmt.Errorf(
				"%s not ready (%s)",
				name,
				state.String(),
			)
		}
	}

	return nil
}
