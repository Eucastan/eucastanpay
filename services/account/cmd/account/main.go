// Package main Account Service API
//
// @title           EucastanPay Account Service API
// @version         1.0
// @description     Authentication and Account Management Service for EucastanPay.
//
// @contact.name    Eucastan
// @contact.email   support@eucastanpay.com
//
// @license.name    MIT
//
// @host account-y0no.onrender.com
// @BasePath /api/v1
// @schemes https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/consumer"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/logger"
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/common/proto/account"
	"github.com/Eucastan/eucastanpay/services/account/config"
	_ "github.com/Eucastan/eucastanpay/services/account/docs"
	"github.com/Eucastan/eucastanpay/services/account/internal/api"
	"github.com/Eucastan/eucastanpay/services/account/internal/api/handler"
	"github.com/Eucastan/eucastanpay/services/account/internal/eventhandler"
	"github.com/Eucastan/eucastanpay/services/account/internal/grpcserver"
	"github.com/Eucastan/eucastanpay/services/account/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/account/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/account/internal/usecase/service"
	"github.com/Eucastan/eucastanpay/services/account/internal/worker"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.SharedCfg.LogLevel)
	log.Info("Starting Account Service...")

	tracer := otel.Tracer("account-service")
	meter := otel.Meter("account-service")

	tm, err := telemetry.New(tracer, meter, log)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize telemetry")
	}

	db := database.NewPostgresDB(cfg, log)
	defer db.CloseDB()

	publisher := producer.NewPublisher(
		cfg.SharedCfg.Kafka.Brokers, cfg.SharedCfg.Kafka.Username,
		cfg.SharedCfg.Kafka.Password, tm,
	)
	defer publisher.Close()

	idempotencyStore := idempotency.NewPostgresStore()
	consumerInit := consumer.NewConsumer(
		cfg.SharedCfg.Kafka.Brokers, cfg.SharedCfg.Kafka.Username,
		cfg.SharedCfg.Kafka.Password, "account-service-group", tm, log,
	)

	accRepo := postgres.NewAccountRepository(db.DB, tm, log)
	accUseCase := service.NewAccountUseCase(accRepo, tm, log)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	go worker.StartOutboxWorker(appCtx, db.DB, publisher, log)

	accountConsumer := eventhandler.NewAccountConsumer(accRepo, accUseCase, idempotencyStore, publisher, tm, log)

	consumerInit.Register(events.TopicUserRegistered,
		consumer.RetryHandler(
			accountConsumer.OnCreateAccountRequest,
			publisher,
			events.TopicUserRegistered,
			events.TopicAccountDLQ,
			tm,
			3,
		),
	)
	consumerInit.Register(events.TopicDebitRequested,
		consumer.RetryHandler(
			accountConsumer.OnDebitRequested,
			publisher,
			events.TopicDebitRequested,
			events.TopicAccountDLQ,
			tm,
			3,
		),
	)
	consumerInit.Register(events.TopicCreditRequested,
		consumer.RetryHandler(
			accountConsumer.OnCreditRequested,
			publisher,
			events.TopicCreditRequested,
			events.TopicAccountDLQ,
			tm,
			3,
		),
	)

	consumerInit.Start(appCtx)

	accHandler := handler.NewAccountHandler(accUseCase)

	// Health check init
	healthChecker := healthcheck.NewHealthChecker("account-service", cfg.Version, log)
	healthChecker.SetDatabase(db.DB)
	healthChecker.SetKafkaProducer(publisher)
	// healthChecker.AddGRPCClient("account-service", allClients.ConnAccount)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", healthChecker.Health)
	r.GET("/live", healthChecker.Liveness)
	r.GET("/ready", healthChecker.Readiness)

	mw := middleware.New(log, cfg.SharedCfg.JWTSecret)
	r.Use(mw.Recovery())
	r.Use(middleware.CorrelationMiddleware())
	r.Use(otelgin.Middleware("account-service"))
	r.Use(mw.Logger(), mw.Auth())

	api.NewRouter(r, accHandler, cfg)

	httpSrv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	// gRPC Server (inter-service)
	listenAddr, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.WithError(err).Fatal("gRPC server failed to listen for connection")
	}

	defer listenAddr.Close()

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(cfg.SharedCfg.JWTSecret)),
	)
	srv := grpcserver.NewAccountServiceServer(accUseCase, accRepo)
	account.RegisterAccountServiceServer(grpcSrv, srv)

	// Start servers
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("HTTP server error")
		}
	}()

	go func() {
		log.Infof("gRPC server listening on :%s", cfg.GRPCPort)
		if err := grpcSrv.Serve(listenAddr); err != nil {
			log.WithError(err).Error("gRPC server error")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Info("Shutting down servers...")

	appCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := consumerInit.Close(); err != nil {
		log.WithError(err).Error("failed to close consumer")
	}

	grpcSrv.GracefulStop()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("failed to shutdown http server")
	}

}
