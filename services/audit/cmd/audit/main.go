package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/consumer"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/logger"
	"github.com/Eucastan/eucastanpay/common/pkg/metrics"
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/audit/config"
	"github.com/Eucastan/eucastanpay/services/audit/internal/api"
	"github.com/Eucastan/eucastanpay/services/audit/internal/api/handler"
	"github.com/Eucastan/eucastanpay/services/audit/internal/eventhandler"
	"github.com/Eucastan/eucastanpay/services/audit/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/audit/internal/infra/tracing"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/audit/internal/usecase/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	log := logger.New(cfg.LogLevel)
	log.Info("Starting Audit Service...")

	db := database.NewPostgresDB(cfg, log)
	defer db.CloseDB()

	publisher := producer.NewPublisher(cfg.Kafka.Brokers)
	defer publisher.Close()

	// Repositories & UseCases
	auditRepo := postgres.NewAuditRepository(db.DB)
	auditUC := service.NewAuditUseCase(auditRepo)

	// tracing
	tracing.InitTracer("audit-service")

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	// Kafka Consumer
	consumerInit := consumer.NewConsumer(cfg.Kafka.Brokers, "audit-service-group", log)
	idempotencyStore := idempotency.NewPostgresStore()

	auditConsumer := eventhandler.NewAuditConsumer(auditRepo, idempotencyStore, log)

	// Register multiple topics
	topics := []string{
		events.TopicTransferInitiated,
		events.TopicTransferCompleted,
		events.TopicTransferFailed,
		events.TopicDebitCompleted,
		events.TopicCreditCompleted,
		events.TopicLedgerCreated,
		events.TopicDebitRequested,
		events.TopicCreditRequested,
	}

	for _, topic := range topics {
		consumerInit.Register(topic, consumer.RetryHandler(
			auditConsumer.Handler(topic),
			publisher,
			topic,
			events.TopicAuditDLQ,
			3,
		))
	}

	consumerInit.Start(appCtx)

	// HTTP Server
	auditHandler := handler.NewAuditHandler(auditUC)
	r := gin.Default()

	// Metrics
	metrics.InitMetrics()
	r.Use(metrics.MetricsMiddleware())

	// Health check init
	healthChecker := healthcheck.NewHealthChecker("audit-service", cfg.Version, log)
	healthChecker.SetDatabase(db.DB)
	healthChecker.SetKafkaProducer(publisher)
	// healthChecker.AddGRPCClient("account-service", allClients.ConnAccount)

	mw := middleware.New(log, cfg.JWTSecret)
	r.Use(mw.Logger(), mw.Recovery(), mw.Auth())
	r.Use(middleware.CorrelationMiddleware())

	r.GET("/health", healthChecker.Health)
	r.GET("/ready", healthChecker.Readiness)
	r.GET("/live", healthChecker.Liveness)

	api.NewRouter(r, auditHandler)

	httpSrv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start HTTP server
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("HTTP server error")
		}
	}()

	log.Infof("Audit Service started on port %s", cfg.HTTPPort)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Info("Shutting down Audit Service...")

	appCancel()

	ShutdownCtx, ShutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ShutdownCancel()

	if err := consumerInit.Close(); err != nil {
		log.WithError(err).Error("Audit consumer shutdown error")
	}

	if err := httpSrv.Shutdown(ShutdownCtx); err != nil {
		log.WithError(err).Error("HTTP shutdown error")
	}

}
