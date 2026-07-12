// Package main Audit Service API
//
// @title           EucastanPay Audit Service API
// @version         1.0
// @description     Audit Service for EucastanPay.
//
// @contact.name    Eucastan
// @contact.email   support@eucastanpay.com
//
// @license.name    MIT
//
// @host eucastanpay.onrender.com
// @BasePath /api/v1
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter: Bearer <JWT>
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
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/audit/config"
	_ "github.com/Eucastan/eucastanpay/services/audit/docs"
	"github.com/Eucastan/eucastanpay/services/audit/internal/api"
	"github.com/Eucastan/eucastanpay/services/audit/internal/api/handler"
	"github.com/Eucastan/eucastanpay/services/audit/internal/eventhandler"
	"github.com/Eucastan/eucastanpay/services/audit/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/audit/internal/usecase/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	log := logger.New(cfg.LogLevel)
	log.Info("Starting Audit Service...")

	tracer := otel.Tracer("audit-service")
	meter := otel.Meter("audit-service")

	tm, err := telemetry.New(tracer, meter, log)
	if err != nil {
		panic(err)
	}

	db := database.NewPostgresDB(cfg, log)
	defer db.CloseDB()

	publisher := producer.NewPublisher(
		cfg.Kafka.Brokers, cfg.Kafka.Username,
		cfg.Kafka.Password, tm,
	)
	defer publisher.Close()

	// Repositories & UseCases
	auditRepo := postgres.NewAuditRepository(db.DB, tm)
	auditUC := service.NewAuditUseCase(auditRepo, tm)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	// Kafka Consumer
	consumerInit := consumer.NewConsumer(
		cfg.Kafka.Brokers, cfg.Kafka.Username,
		cfg.Kafka.Password, "audit-service-group", tm, log,
	)
	idempotencyStore := idempotency.NewPostgresStore()

	auditConsumer := eventhandler.NewAuditConsumer(auditRepo, idempotencyStore, tm, log)

	// Register multiple topics
	topics := []string{
		events.TopicUserRegistered,
		events.TopicUserRegistrationFailed,
		events.TopicUserKYCCreated,
		events.TopicUserKYCVerified,
		events.TopicAccountCreated,
		events.TopicCreateAccFailed,
		events.TopicDepositAccount,
		events.TopicWithdrawal,
		events.TopicTransferInitiated,
		events.TopicReverseInitiated,
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
			tm,
			3,
		))
	}

	consumerInit.Start(appCtx)

	// HTTP Server
	auditHandler := handler.NewAuditHandler(auditUC)
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check init
	healthChecker := healthcheck.NewHealthChecker("audit-service", cfg.Version, log)
	healthChecker.SetDatabase(db.DB)
	healthChecker.SetKafkaProducer(publisher)
	// healthChecker.AddGRPCClient("account-service", allClients.ConnAccount)

	mw := middleware.New(log, cfg.JWTSecret)
	r.Use(mw.Recovery())
	r.Use(middleware.CorrelationMiddleware())
	r.Use(otelgin.Middleware("audit-service"))
	r.Use(mw.Logger())

	r.GET("/health", healthChecker.Health)
	r.GET("/ready", healthChecker.Readiness)
	r.GET("/live", healthChecker.Liveness)

	api.NewRouter(r, auditHandler, cfg)

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
