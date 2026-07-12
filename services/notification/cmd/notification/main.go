// Package main Notification Service API
//
// @title           EucastanPay Notification Service API
// @version         1.0
// @description     Notification Management Service for EucastanPay.
//
// @contact.name    Eucastan
// @contact.email   support@eucastanpay.com
//
// @license.name    MIT
//
// @host eucastanpay.onrender.com
// @BasePath /api/v1
// @schemes https http
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
	"github.com/Eucastan/eucastanpay/services/notification/config"
	_ "github.com/Eucastan/eucastanpay/services/notification/docs"
	"github.com/Eucastan/eucastanpay/services/notification/internal/api"
	"github.com/Eucastan/eucastanpay/services/notification/internal/api/handler"
	"github.com/Eucastan/eucastanpay/services/notification/internal/eventhandler"
	"github.com/Eucastan/eucastanpay/services/notification/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/notification/internal/provider"
	"github.com/Eucastan/eucastanpay/services/notification/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/notification/internal/usecase/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.SharedCfg.LogLevel)
	log.Info("Starting Notification Service...")

	tracer := otel.Tracer("notification-service")
	meter := otel.Meter("notification-service")

	tm, err := telemetry.New(tracer, meter, log)
	if err != nil {
		panic(err)
	}

	db := database.NewPostgresDB(cfg, log)
	defer db.CloseDB()

	emailProvider := provider.NewEmailProvider(cfg)

	publisher := producer.NewPublisher(
		cfg.SharedCfg.Kafka.Brokers, cfg.SharedCfg.Kafka.Username,
		cfg.SharedCfg.Kafka.Password, tm,
	)

	consumerInit := consumer.NewConsumer(
		cfg.SharedCfg.Kafka.Brokers, cfg.SharedCfg.Kafka.Username,
		cfg.SharedCfg.Kafka.Password, "notification-service-group",
		tm, log,
	)

	notificationRepo := postgres.NewNotificationRepository(db.DB, tm)
	notificationUC := service.NewNotificationUseCase(notificationRepo, tm, emailProvider, log)

	idempotencyStore := idempotency.NewPostgresStore()
	notificationConsumer := eventhandler.NewNotificationConsumer(
		notificationRepo, idempotencyStore, tm, log,
	)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	consumerInit.Register(events.TopicUserRegistered,
		consumer.RetryHandler(
			notificationConsumer.OnUserRegistered,
			publisher,
			events.TopicUserRegistered,
			events.TopicNotificationDLQ,
			tm,
			3,
		),
	)

	consumerInit.Register(events.TopicUserKYCCreated,
		consumer.RetryHandler(
			notificationConsumer.OnKycCreated,
			publisher,
			events.TopicUserKYCCreated,
			events.TopicNotificationDLQ,
			tm,
			3,
		),
	)

	consumerInit.Register(events.TopicUserKYCVerified,
		consumer.RetryHandler(
			notificationConsumer.OnUserKycVerified,
			publisher,
			events.TopicUserKYCVerified,
			events.TopicNotificationDLQ,
			tm,
			3,
		),
	)

	consumerInit.Register(events.TopicAccountCreated,
		consumer.RetryHandler(
			notificationConsumer.OnAccountCreated,
			publisher,
			events.TopicAccountCreated,
			events.TopicNotificationDLQ,
			tm,
			3,
		),
	)

	consumerInit.Register(events.TopicCreateAccFailed,
		consumer.RetryHandler(
			notificationConsumer.OnAccountCreationFailed,
			publisher,
			events.TopicCreateAccFailed,
			events.TopicNotificationDLQ,
			tm,
			3,
		),
	)

	consumerInit.Register(events.TopicDepositAccount,
		consumer.RetryHandler(
			notificationConsumer.OnAccountDeposit,
			publisher,
			events.TopicDepositAccount,
			events.TopicNotificationDLQ,
			tm,
			3,
		),
	)

	consumerInit.Register(events.TopicWithdrawal,
		consumer.RetryHandler(
			notificationConsumer.OnCashWithdraw,
			publisher,
			events.TopicWithdrawal,
			events.TopicNotificationDLQ,
			tm,
			3,
		),
	)

	consumerInit.Register(events.TopicTransferCompleted,
		consumer.RetryHandler(
			notificationConsumer.OnTransferCompleted,
			publisher,
			events.TopicTransferCompleted,
			events.TopicNotificationDLQ,
			tm,
			3,
		),
	)

	consumerInit.Register(events.TopicTransferFailed,
		consumer.RetryHandler(
			notificationConsumer.OnTransferFailed,
			publisher,
			events.TopicTransferFailed,
			events.TopicNotificationDLQ,
			tm,
			3,
		),
	)

	consumerInit.Start(appCtx)

	notificationHandler := handler.NewNotificationHandler(notificationUC)

	// Health check init
	healthChecker := healthcheck.NewHealthChecker("notification-service", cfg.Version, log)
	healthChecker.SetDatabase(db.DB)
	healthChecker.SetKafkaConsumer(consumerInit)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", healthChecker.Health)
	r.GET("/live", healthChecker.Liveness)
	r.GET("/ready", healthChecker.Readiness)

	mw := middleware.New(log, cfg.SharedCfg.JWTSecret)
	r.Use(mw.Recovery())
	r.Use(middleware.CorrelationMiddleware())
	r.Use(otelgin.Middleware("notification-service"))
	r.Use(mw.Logger(), mw.Auth())

	api.NewRouter(r, notificationHandler)

	httpSrv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("HTTP server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Info("Shutting down Notification service")

	appCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := consumerInit.Close(); err != nil {
		log.WithError(err).Error("Notification service consumer shutdown error")
	}

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("Notification service shutdown error")
	}
}
