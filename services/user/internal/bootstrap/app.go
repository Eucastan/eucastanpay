package bootstrap

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/user/config"
	"github.com/Eucastan/eucastanpay/services/user/internal/grpcserver"
	"github.com/Eucastan/eucastanpay/services/user/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/user/internal/infra/redis"
	"github.com/Eucastan/eucastanpay/services/user/internal/repository"
	"github.com/Eucastan/eucastanpay/services/user/internal/usecase"
	"github.com/Eucastan/eucastanpay/services/user/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type App struct {
	server       *http.Server
	router       *gin.Engine
	grpcServ     *grpc.Server
	cfg          *config.Config
	logger       *logrus.Logger
	telemetry    *telemetry.Telemetry
	health       *healthcheck.Checker
	publish      *producer.Publisher
	redis        *redis.RedisClient
	worker       *worker.OutboxWorker
	workerCtx    context.Context
	workerCancel context.CancelFunc
	database     *database.DBConnect
	userRepo     repository.UserRepository
	authRepo     repository.AuthRepository
	kycRepo      repository.KYCRepository
	userUC       usecase.UserUseCaseInterface
	kycUC        usecase.KYCUseCase
	email        usecase.EmailSender
	grpcserver   *grpcserver.UserServer
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	app := &App{
		cfg: cfg,
	}

	if err := app.bootstrap(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run() error {
	errCh := make(chan error, 2)

	go func() {
		a.logger.Infof(
			"User service listening on %s",
			a.cfg.HTTPPort,
		)

		errCh <- a.server.ListenAndServe()
	}()

	go func() {
		grpcListener, err := net.Listen("tcp", ":"+a.cfg.GRPCPort)
		if err != nil {
			errCh <- err
			return
		}

		a.logger.Infof("User gRPC server listening on %s", a.cfg.GRPCPort)

		errCh <- a.grpcServ.Serve(grpcListener)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}

	case <-stop:
		a.logger.Info("shutdown signal received")
		return a.shutdown(context.Background())
	}

	return nil
}

func (a *App) bootstrap() error {
	a.initLogger()

	if err := a.initTelemetry(); err != nil {
		return err
	}

	a.initDatabase()
	a.initEmail()
	a.initRedis()
	a.initPublisher()
	a.initRepository()
	a.initUseCase()
	a.initHealth()
	a.initRouter()

	a.initWorkerContext()
	a.initOutboxWorker()

	go a.worker.Start(a.workerCtx)

	a.initGRPCServer()
	a.initServer()

	return nil
}
