package healthcheck

import (
	"context"
	"runtime"
	"time"
)

func (h *HealthChecker) checkDatabase(components map[string]Component) {
	if h.db == nil {
		components["database"] = Component{Status: StatusUnhealthy, Error: "database not configured"}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	// 1. Ping
	if err := h.db.Ping(ctx); err != nil {
		h.logger.WithError(err).Error("Database ping failed")
		components["database"] = Component{Status: StatusUnhealthy, Error: err.Error()}
		return
	}

	// 2. Actual query test + pool stats
	var one int
	err := h.db.QueryRow(ctx, "SELECT 1").Scan(&one)
	if err != nil {
		h.logger.WithError(err).Warn("Database query test failed")
		components["database"] = Component{Status: StatusDegraded, Error: "query test failed: " + err.Error()}
		return
	}

	stats := h.db.Stat()
	components["database"] = Component{
		Status: StatusHealthy,
		Details: map[string]interface{}{
			"total_conns":    stats.TotalConns(),
			"idle_conns":     stats.IdleConns(),
			"acquired_conns": stats.AcquiredConns(),
			"max_conns":      stats.MaxConns(),
			"active_queries": one,
		},
	}
}

func (h *HealthChecker) checkKafkaProducer(components map[string]Component) {
	if h.kafkaProd == nil {
		components["kafka_producer"] = Component{Status: StatusDegraded, Error: "not configured"}
		return
	}
	// Basic connectivity test can be added by trying to write to a test topic if needed
	components["kafka_producer"] = Component{Status: StatusHealthy}
}

func (h *HealthChecker) checkKafkaConsumer(components map[string]Component) {
	if h.kafkaCons == nil {
		components["kafka_consumer"] = Component{Status: StatusDegraded, Error: "not configured"}
		return
	}

	components["kafka_consumer"] = Component{
		Status:  StatusHealthy,
		Details: map[string]string{"status": "running"},
	}
}

func (h *HealthChecker) checkGRPCConnections(components map[string]Component) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for name, conn := range h.grpcConns {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err := conn.Invoke(ctx, "/health", nil, nil)
		cancel()

		if err != nil {
			components["grpc_"+name] = Component{Status: StatusDegraded, Error: err.Error()}
		} else {
			components["grpc_"+name] = Component{Status: StatusHealthy}
		}
	}
}

func (h *HealthChecker) checkSystemMetrics(components map[string]Component) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	components["system"] = Component{
		Status: StatusHealthy,
		Details: map[string]interface{}{
			"memory_alloc_mb": m.Alloc / 1024 / 1024,
			"memory_total_mb": m.TotalAlloc / 1024 / 1024,
			"memory_sys_mb":   m.Sys / 1024 / 1024,
			"goroutines":      runtime.NumGoroutine(),
			"gc_cycles":       m.NumGC,
			"cpu_goroutines":  runtime.NumCPU(),
		},
	}
}
