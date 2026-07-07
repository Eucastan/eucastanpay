// Package main User Service API
//
// @title           EucastanPay User Service API
// @version         1.0
// @description     Authentication and User Management Service for EucastanPay.
//
// @contact.name    Eucastan
// @contact.email   support@eucastanpay.com
//
// @license.name    MIT
//
// @host localhost:8001
// @BasePath /api/v1
// @schemes http https
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

	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/logger"
	"github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/common/proto/user"
	"github.com/Eucastan/eucastanpay/services/user/config"
	_ "github.com/Eucastan/eucastanpay/services/user/docs"
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
	log.Info("Starting User service...")

	tracer := otel.Tracer("user-service")
	meter := otel.Meter("user-service")

	tm, err := telemetry.New(tracer, meter, log)
	if err != nil {
		panic(err)
	}

	db := database.NewPostgresDB(cfg, log)
	defer db.CloseDB()

	email := usecase.NewEmailService(cfg)

	redis := redis.NewRedisClient(cfg, log)
	defer redis.Close()

	publisher := producer.NewPublisher(cfg.Kafka.Brokers, tm)
	defer publisher.Close()

	authRepo := postgres.NewAuthRepository(db.DB, tm)
	kycRepo := postgres.NewKYCRepository(db.DB, tm)
	userRepo := postgres.NewUserRepository(db.DB, tm)

	kycUseCase := service.NewKYCUseCase(kycRepo, tm, cfg)

	registerEvent := worker.NewPublishUserRegistration(userRepo)
	userUseCase := service.NewUserUseCase(
		userRepo, authRepo, tm, cfg, email,
		redis, registerEvent,
	)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	go worker.StartOutboxWorker(appCtx, db.DB, publisher, tm, log)

	kycHandler := handler.NewKYCHandler(kycUseCase)
	userHandler := handler.NewUserHandler(userUseCase)

	// Health check init
	healthChecker := healthcheck.NewHealthChecker("user-service", cfg.Version, log)
	healthChecker.SetDatabase(db.DB)
	healthChecker.SetKafkaProducer(publisher)
	// healthChecker.AddGRPCClient("account-service", allClients.ConnAccount)

	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	mw := middleware.New(log, cfg.JWTSecret)
	r.Use(mw.Recovery(), mw.Logger())
	r.Use(middleware.CorrelationMiddleware())
	r.Use(otelgin.Middleware("notification-service"))

	r.GET("/health", healthChecker.Health)
	r.GET("/live", healthChecker.Liveness)
	r.GET("/ready", healthChecker.Readiness)

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

	appCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("User service shutdown error")
	}

	grpcSrv.GracefulStop()
}
