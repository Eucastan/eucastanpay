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
	"github.com/Eucastan/eucastanpay/common/proto/account"
	"github.com/Eucastan/eucastanpay/services/account/config"
	"github.com/Eucastan/eucastanpay/services/account/internal/api"
	"github.com/Eucastan/eucastanpay/services/account/internal/api/handler"
	"github.com/Eucastan/eucastanpay/services/account/internal/eventhandler"
	"github.com/Eucastan/eucastanpay/services/account/internal/grpcserver"
	"github.com/Eucastan/eucastanpay/services/account/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/account/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/account/internal/usecase/service"
	"github.com/Eucastan/eucastanpay/services/account/internal/worker"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)
	log.Info("Starting Account Service...")

	db := database.NewPostgresDB(cfg, log)
	defer db.CloseDB()

	publisher := producer.NewPublisher(cfg.Kafka.Brokers)
	defer publisher.Close()

	idempotencyStore := idempotency.NewPostgresStore()
	consumerInit := consumer.NewConsumer(cfg.Kafka.Brokers, "account-service-group", log)

	accRepo := postgres.NewAccountRepository(db.DB, log)
	accUseCase := service.NewAccountUseCase(accRepo, log)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	go worker.StartOutboxWorker(appCtx, db.DB, publisher, log)

	accountConsumer := eventhandler.NewAccountConsumer(accRepo, accUseCase, idempotencyStore, publisher, log)

	consumerInit.Register(events.TopicUserRegistered,
		consumer.RetryHandler(
			accountConsumer.OnUserRegistration,
			publisher,
			events.TopicUserRegistered,
			events.TopicAccountDLQ,
			3,
		),
	)
	consumerInit.Register(events.TopicDebitRequested,
		consumer.RetryHandler(
			accountConsumer.OnDebitRequested,
			publisher,
			events.TopicDebitRequested,
			events.TopicAccountDLQ,
			3,
		),
	)
	consumerInit.Register(events.TopicCreditRequested,
		consumer.RetryHandler(
			accountConsumer.OnCreditRequested,
			publisher,
			events.TopicCreditRequested,
			events.TopicAccountDLQ,
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
	mw := middleware.New(log, cfg.JWTSecret)
	r.Use(mw.Logger(), mw.Recovery())
	r.Use(middleware.CorrelationMiddleware())

	r.GET("/health", healthChecker.Health)
	r.GET("/live", healthChecker.Liveness)
	r.GET("/ready", healthChecker.Readiness)

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
		grpc.UnaryInterceptor(interceptor.AuthInterceptor(cfg.JWTSecret)),
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
