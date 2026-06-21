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
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/consumer"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/logger"
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/transfer/config"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/api"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/api/handler"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/eventhandler"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/grpc/clients"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/usecase/service"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/worker"
	"github.com/gin-gonic/gin"
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

	go worker.StartOutboxWorker(context.Background(), db.DB, publisher, log)
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

	consumerInit.Start(context.Background())

	// Initialize gRPC Clients
	allClients, err := clients.Init(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize gRPC clients: ", err)
	}
	defer allClients.Close()

	transferUC := service.NewTransferUseCase(transferRepo, allClients, publisher, log)
	transferHandler := handler.NewTransferHandler(transferUC)
	go transferUC.RecoverStuckTransfers(context.Background())

	r := gin.Default()
	mw := middleware.New(log, cfg.JWTSecret)
	r.Use(mw.Logger(), mw.Recovery(), mw.Auth())
	r.Use(middleware.CorrelationMiddleware())
	api.NewRouter(r, transferHandler)

	httpSrv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Graceful Shutdown
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("HTTP server error")
		}
	}()

	log.Infof("Transfer Service started successfully on port %s", cfg.HTTPPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Info("Shutting down Transfer Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpSrv.Shutdown(ctx)
}
