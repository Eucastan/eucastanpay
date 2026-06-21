package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/logger"
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/common/proto/user"
	"github.com/Eucastan/eucastanpay/services/user/config"
	"github.com/Eucastan/eucastanpay/services/user/internal/api"
	"github.com/Eucastan/eucastanpay/services/user/internal/api/handler"
	"github.com/Eucastan/eucastanpay/services/user/internal/grpcserver"
	"github.com/Eucastan/eucastanpay/services/user/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/user/internal/infra/redis"
	"github.com/Eucastan/eucastanpay/services/user/internal/repository/postgres"
	"github.com/Eucastan/eucastanpay/services/user/internal/usecase"
	"github.com/Eucastan/eucastanpay/services/user/internal/usecase/service"
	"github.com/Eucastan/eucastanpay/services/user/internal/worker"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)
	log.Info("Starting User service...")

	db := database.NewPostgresDB(cfg, log)
	defer db.CloseDB()

	email := usecase.NewEmailService(cfg)

	redis := redis.NewRedisClient(cfg, log)
	defer redis.Close()

	publisher := producer.NewPublisher(cfg.Kafka.Brokers)
	defer publisher.Close()

	authRepo := postgres.NewAuthRepository(db.DB)
	kycRepo := postgres.NewKYCRepository(db.DB)
	userRepo := postgres.NewUserRepository(db.DB)

	kycUseCase := service.NewKYCUseCase(kycRepo, cfg)

	registerEvent := worker.NewPublishUserRegistration(userRepo)
	userUseCase := service.NewUserUseCase(
		userRepo, authRepo, cfg, email,
		redis, registerEvent,
	)

	go worker.StartOutboxWorker(context.Background(), db.DB, publisher, log)

	kycHandler := handler.NewKYCHandler(kycUseCase)
	userHandler := handler.NewUserHandler(userUseCase)

	r := gin.New()
	mw := middleware.New(log, cfg.JWTSecret)
	r.Use(mw.Recovery(), mw.Logger())
	r.Use(middleware.CorrelationMiddleware())

	api.NewRouter(r, userHandler, kycHandler, cfg)

	httpSrv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// gRPC Server (inter-service)
	listenAddr, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Info(err)
	}

	defer listenAddr.Close()

	grpcSrv := grpc.NewServer()
	srv := grpcserver.NewUserServiceServer(userUseCase, kycUseCase)
	user.RegisterUserServiceServer(grpcSrv, srv)

	// Graceful shutdown
	go httpSrv.ListenAndServe()
	go grpcSrv.Serve(listenAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	log.Info("Shutting Down User Service")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpSrv.Shutdown(ctx)
	grpcSrv.GracefulStop()
}
