package bootstrap

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func (a *App) shutdown(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, a.cfg.ShutdownTimeout)
	defer cancel()

	a.logger.Info("Starting graceful shutdown...")

	// Stop receiving new traffic
	if a.health != nil {
		a.health.SetReady(false)
	}

	// Give load balancer time to stop routing traffic.
	time.Sleep(3 * time.Second)

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			a.logger.Errorf("HTTP shutdown failed: %v", err)
		}
	}

	a.logger.Info("HTTP server stopped")

	if a.manager != nil {
		if err := a.manager.Close(); err != nil {
			a.logger.Errorf("Failed closing gRPC connections: %v", err)
		}
	}

	a.logger.Info("gRPC connections closed")

	/*if a.telemetry != nil {
		if err := a.telemetry.Shutdown(ctx); err != nil {
			a.logger.Errorf("Telemetry shutdown failed: %v", err)
		}
	}

	a.logger.Info("Telemetry stopped")*/
	a.logger.Info("Graceful shutdown completed")

	return nil
}
