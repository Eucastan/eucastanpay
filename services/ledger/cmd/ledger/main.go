// Package main Ledger Service API
//
// @title           EucastanPay Ledger Service API
// @version         1.0
// @description     Authentication and Ledger Entry Service for EucastanPay.
//
// @contact.name    Eucastan
// @contact.email   support@eucastanpay.com
//
// @license.name    MIT
//
// @host localhost:8003
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
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Eucastan/eucastanpay/common/idempotency"
	"github.com/Eucastan/eucastanpay/common/pkg/events"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/clients"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/interceptor"
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/consumer"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/logger"
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/common/proto/ledger"
	"github.com/Eucastan/eucastanpay/services/ledger/config"
	_ "github.com/Eucastan/eucastanpay/services/ledger/docs"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/api"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/api/handler"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/eventshandler"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/grpcserver"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/usecase/service"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/worker"
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

	log := logger.New(cfg.LogLevel)
	log.Info("Starting Ledger Service...")

	tracer := otel.Tracer("ledger-service")
	meter := otel.Meter("ledger-service")

	tm, err := telemetry.New(tracer, meter, log)
	if err != nil {
		panic(err)
	}

	db := database.NewPostgresDB(cfg, log)
	defer db.CloseDB()

	publisher := producer.NewPublisher(cfg.Kafka.Brokers, tm)
	defer publisher.Close()

	accountConfig := clients.Config{
		AccountServiceAddr: cfg.AccountGRPCPort,
		Timeout:            5 * time.Second,
		MaxRetries:         3,
		Insecure:           true,
	}
	clients, err := clients.NewClients(accountConfig, log)
	if err != nil {
		log.WithError(err).Fatal("failed to connect gRPC clients")
	}

	ledgerRepo := postgres.NewLedgerRepository(db.DB, tm)
	ledgerUC := service.NewLedgerUseCase(ledgerRepo, tm, clients, log)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	idemStore := idempotency.NewPostgresStore()
	go worker.StartOutboxWorker(appCtx, db.DB, publisher, log)

	consumerInit := consumer.NewConsumer(cfg.Kafka.Brokers, "ledger-group", tm, log)
	ledgerConsumer := eventshandler.NewLedgerEventHandler(ledgerRepo, ledgerUC, tm, idemStore, publisher, log)

	consumerInit.Register(events.TopicTransferCompleted,
		consumer.RetryHandler(
			ledgerConsumer.OnTransferCompleted,
			publisher,
			events.TopicTransferCompleted,
			events.TopicLedgerDLQ,
			tm,
			3,
		),
	)

	consumerInit.Register(events.TopicDepositAccount,
		consumer.RetryHandler(
			ledgerConsumer.OnAccountDeposit,
			publisher,
			events.TopicDepositAccount,
			events.TopicLedgerDLQ,
			tm,
			3,
		),
	)

	consumerInit.Register(events.TopicWithdrawal,
		consumer.RetryHandler(
			ledgerConsumer.OnCasWithdraw,
			publisher,
			events.TopicWithdrawal,
			events.TopicLedgerDLQ,
			tm,
			3,
		),
	)

	consumerInit.Start(appCtx)
	ledgerHandler := handler.NewLedgerHandler(ledgerUC)

	// Health check init
	healthChecker := healthcheck.NewHealthChecker("ledger-service", cfg.Version, log)
	healthChecker.SetDatabase(db.DB)
	healthChecker.SetKafkaProducer(publisher)
	healthChecker.AddGRPCClient("account-service", clients.ConnAccount)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	mw := middleware.New(log, cfg.JWTSecret)
	r.Use(mw.Recovery())
	r.Use(middleware.CorrelationMiddleware())
	r.Use(otelgin.Middleware("ledger-service"))
	r.Use(mw.Logger())

	r.GET("/health", healthChecker.Health)
	r.GET("/live", healthChecker.Liveness)
	r.GET("/ready", healthChecker.Readiness)

	api.NewRouter(r, ledgerHandler, cfg)

	httpSrv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	listenAddr, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.WithError(err).Error("failed to connect to gRPC server")
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(cfg.JWTSecret)),
	)
	srv := grpcserver.NewLedgerServiceServer(ledgerUC)
	ledger.RegisterLedgerServiceServer(grpcServer, srv)

	// Start servers
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("HTTP server error")
		}
	}()

	go func() {
		log.Infof("gRPC server listening on %s", cfg.GRPCPort)
		if err := grpcServer.Serve(listenAddr); err != nil {
			log.WithError(err).Error("gRPC server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Info("Shutting down Ledger Service...")

	appCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := consumerInit.Close(); err != nil {
		log.WithError(err).Error("Ledger service consumer error")
	}

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("Ledger service shutdown error")
	}

	grpcServer.GracefulStop()
}
