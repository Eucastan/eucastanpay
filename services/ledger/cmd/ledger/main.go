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
	"github.com/Eucastan/eucastanpay/common/proto/ledger"
	"github.com/Eucastan/eucastanpay/services/ledger/config"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/api"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/api/handler"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/eventshandler"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/grpcserver"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/usecase/service"
	"github.com/Eucastan/eucastanpay/services/ledger/internal/worker"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)
	log.Info("Starting Ledger Service...")

	db := database.NewPostgresDB(cfg, log)
	defer db.CloseDB()

	publisher := producer.NewPublisher(cfg.Kafka.Brokers)
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

	ledgerRepo := postgres.NewLedgerRepository(db.DB)
	ledgerUC := service.NewLedgerUseCase(ledgerRepo, clients, log)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	idemStore := idempotency.NewPostgresStore()
	go worker.StartOutboxWorker(appCtx, db.DB, publisher, log)

	consumerInit := consumer.NewConsumer(cfg.Kafka.Brokers, "ledger-group", log)
	ledgerConsumer := eventshandler.NewLedgerEventHandler(ledgerRepo, ledgerUC, idemStore, publisher, log)

	consumerInit.Register(events.TopicTransferCompleted,
		consumer.RetryHandler(
			ledgerConsumer.OnTransferCompleted,
			publisher,
			events.TopicTransferCompleted,
			events.TopicLedgerDLQ,
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
	mw := middleware.New(log, cfg.JWTSecret)
	r.Use(mw.Logger(), mw.Recovery())
	r.Use(middleware.CorrelationMiddleware())

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
