package healthcheck

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/kafka/consumer"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type HealthChecker struct {
	serviceName string
	version     string
	startTime   time.Time
	db          *pgxpool.Pool
	kafkaProd   *producer.Publisher
	kafkaCons   *consumer.Consumer
	grpcConns   map[string]*grpc.ClientConn
	mu          sync.RWMutex
	logger      *logrus.Logger
	panicCount  int64
	failedTx    int64
}

func NewHealthChecker(serviceName, version string, logger *logrus.Logger) *HealthChecker {
	return &HealthChecker{
		serviceName: serviceName,
		version:     version,
		startTime:   time.Now().UTC(),
		grpcConns:   make(map[string]*grpc.ClientConn),
		logger:      logger,
	}
}

// Add dependencies
func (h *HealthChecker) SetDatabase(db *pgxpool.Pool)           { h.db = db }
func (h *HealthChecker) SetKafkaProducer(p *producer.Publisher) { h.kafkaProd = p }
func (h *HealthChecker) SetKafkaConsumer(c *consumer.Consumer)  { h.kafkaCons = c }
func (h *HealthChecker) AddGRPCClient(name string, conn *grpc.ClientConn) {
	h.mu.Lock()
	h.grpcConns[name] = conn
	h.mu.Unlock()
}

func (h *HealthChecker) Health(c *gin.Context) {
	report := HealthReport{
		Status:     "healthy",
		Service:    h.serviceName,
		Version:    h.version,
		Uptime:     time.Since(h.startTime).String(),
		Timestamp:  time.Now().UTC(),
		Components: make(map[string]Component),
		Summary:    make(map[string]interface{}),
	}

	// Run all checks
	h.checkDatabase(report.Components)
	h.checkKafkaProducer(report.Components)
	h.checkKafkaConsumer(report.Components)
	h.checkGRPCConnections(report.Components)
	h.checkSystemMetrics(report.Components)

	// Overall status
	for _, comp := range report.Components {
		if comp.Status == StatusUnhealthy {
			report.Status = "unhealthy"
			break
		}
		if comp.Status == StatusDegraded && report.Status == "healthy" {
			report.Status = "degraded"
		}
	}

	c.JSON(http.StatusOK, report)
}

func (h *HealthChecker) Readiness(c *gin.Context) {
	report := map[string]interface{}{
		"status":    "ready",
		"service":   h.serviceName,
		"version":   h.version,
		"timestamp": time.Now().UTC(),
	}

	components := make(map[string]Component)

	// Critical checks for readiness
	h.checkDatabase(components)
	h.checkKafkaProducer(components)

	ready := true
	for _, comp := range components {
		if comp.Status == StatusUnhealthy {
			ready = false
			break
		}
	}

	if !ready {
		report["status"] = "not ready"
		c.JSON(http.StatusServiceUnavailable, report)
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *HealthChecker) Liveness(c *gin.Context) {
	report := map[string]interface{}{
		"status":     "alive",
		"service":    h.serviceName,
		"version":    h.version,
		"goroutines": runtime.NumGoroutine(),
		"uptime":     time.Since(h.startTime).String(),
		"timestamp":  time.Now().UTC(),
	}

	// Check for too many goroutines (potential leak)
	if runtime.NumGoroutine() > 10000 {
		report["status"] = "warning"
		report["issue"] = "too many goroutines"
	}

	c.JSON(http.StatusOK, report)
}
