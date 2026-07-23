package checks

import (
	"context"
	"runtime"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
)

type System struct{}

func NewSystem() *System {
	return &System{}
}

func (s *System) Name() string {
	return "system"
}

func (s *System) Check(ctx context.Context) healthcheck.Component {

	started := time.Now()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return healthcheck.Component{
		Name:     s.Name(),
		Status:   healthcheck.Healthy,
		Duration: time.Since(started).String(),
		Details: map[string]interface{}{
			"go_version":      runtime.Version(),
			"cpu_count":       runtime.NumCPU(),
			"goroutines":      runtime.NumGoroutine(),
			"memory_alloc_mb": mem.Alloc / 1024 / 1024,
			"memory_sys_mb":   mem.Sys / 1024 / 1024,
			"heap_objects":    mem.HeapObjects,
			"gc_cycles":       mem.NumGC,
		},
	}
}
