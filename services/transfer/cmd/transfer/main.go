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
	"github.com/Eucastan/eucastanpay/common/proto/transfer"
	"github.com/Eucastan/eucastanpay/services/transfer/config"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/api"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/api/handler"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/eventhandler"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/grpc/clients"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/grpc/server"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/usecase/service"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/worker"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)
	log.Info("Starting Transfer Service...")

	db := database.NewPostgresDB(cfg, log)
	defer db.CloseDB()

	transferRepo := postgres.NewTransferRepository(db.DB)
	//redis := redis.NewRedisClient(cfg)

	// Kafka init
	publisher := producer.NewPublisher(cfg.Kafka.Brokers)
	defer publisher.Close()

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	go worker.StartOutboxWorker(appCtx, db.DB, publisher, log)
	idempotencyStore := idempotency.NewPostgresStore()

	consumerInit := consumer.NewConsumer(cfg.Kafka.Brokers, "transfer-group", log)
	transferConsumer := eventhandler.NewTransferConsumer(transferRepo, idempotencyStore, publisher, log)

	consumerInit.Register(events.TopicTransferInitiated,
		consumer.RetryHandler(
			transferConsumer.OnTransferInitiated,
			publisher,
			events.TopicTransferInitiated,
			events.TopicTransferDLQ,
			3,
		),
	)

	consumerInit.Register(events.TopicDebitCompleted,
		consumer.RetryHandler(
			transferConsumer.OnDebitCompleted,
			publisher,
			events.TopicDebitCompleted,
			events.TopicTransferDLQ,
			3,
		),
	)

	consumerInit.Register(events.TopicDebitFailed,
		consumer.RetryHandler(
			transferConsumer.OnDebitFailed,
			publisher,
			events.TopicDebitFailed,
			events.TopicTransferDLQ,
			3,
		),
	)

	consumerInit.Register(events.TopicCreditCompleted,
		consumer.RetryHandler(
			transferConsumer.OnCreditCompleted,
			publisher,
			events.TopicCreditCompleted,
			events.TopicTransferDLQ,
			3,
		),
	)

	consumerInit.Register(events.TopicCreditFailed,
		consumer.RetryHandler(
			transferConsumer.OnCreditFailed,
			publisher,
			events.TopicCreditFailed,
			events.TopicTransferDLQ,
			3,
		),
	)

	consumerInit.Register(
		events.TopicDebitReverseCompleted,
		consumer.RetryHandler(
			transferConsumer.OnDebitReverseCompleted,
			publisher,
			events.TopicDebitReverseCompleted,
			events.TopicTransferDLQ,
			3,
		),
	)

	consumerInit.Start(appCtx)

	// Initialize gRPC Clients
	allClients, err := clients.Init(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize gRPC clients: ", err)
	}
	defer allClients.Close()

	transferUC := service.NewTransferUseCase(transferRepo, allClients, publisher, log)
	transferHandler := handler.NewTransferHandler(transferUC)
	go transferUC.RecoverStuckTransfers(context.Background())

	// Health check init
	healthChecker := healthcheck.NewHealthChecker("transfer-service", cfg.Version, log)
	healthChecker.SetDatabase(db.DB)
	healthChecker.SetKafkaProducer(publisher)
	healthChecker.AddGRPCClient("account-service", allClients.ConnAccount)

	r := gin.Default()
	mw := middleware.New(log, cfg.JWTSecret)
	r.Use(mw.Logger(), mw.Recovery(), mw.Auth())
	r.Use(middleware.CorrelationMiddleware())

	r.GET("/health", healthChecker.Health)
	r.GET("/live", healthChecker.Liveness)
	r.GET("/ready", healthChecker.Readiness)

	api.NewRouter(r, transferHandler)

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
	srv := server.NewTransferServiceServer(transferUC)
	transfer.RegisterTransferServiceServer(grpcServer, srv)

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("HTTP server error")
		}
	}()

	log.Infof("Transfer Service started successfully on port %s", cfg.HTTPPort)

	go func() {
		log.Infof("gRPC server listening on %s", cfg.GRPCPort)
		if err := grpcServer.Serve(listenAddr); err != nil {
			log.WithError(err).Error("gRPC server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Info("Shutting down Transfer Service...")

	appCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := consumerInit.Close(); err != nil {
		log.WithError(err).Error("Transfer service consumer shutdown error")
	}

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("Transfer service shutdown error")
	}

	grpcServer.GracefulStop()
}
