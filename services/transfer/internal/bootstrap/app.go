package bootstrap

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	manager "github.com/Eucastan/eucastanpay/common/pkg/grpc"
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/consumer"
	"github.com/Eucastan/eucastanpay/common/pkg/kafka/producer"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/transfer/config"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/grpc/server"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/infra/database"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/repository"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/usecase"
	"github.com/Eucastan/eucastanpay/services/transfer/internal/worker"
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
	consumer     *consumer.Consumer
	worker       *worker.OutboxWorker
	workerCtx    context.Context
	workerCancel context.CancelFunc
	database     *database.DBConnect
	repo         repository.TransferRepository
	uc           usecase.TransferUseCase
	manager      *manager.Manager
	grpcserver   *server.TransferServiceServer
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
			"Transfer service listening on %s",
			a.cfg.HTTPPort,
		)

		errCh <- a.server.ListenAndServe()
	}()

	go func() {
		grpcListener, err := net.Listen("tcp", ":"+a.cfg.GRPCADDR)
		if err != nil {
			errCh <- err
			return
		}

		a.logger.Infof("Transfer gRPC server listening on %s", a.cfg.GRPCADDR)

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
	a.initPublisher()
	a.initConsumer()

	if err := a.initManager(); err != nil {
		return err
	}

	a.initRepository()
	a.initUseCase()
	a.initHealth()
	a.initRouter()

	a.initWorkerContext()
	a.initOutboxWorker()

	go a.worker.Start(a.workerCtx)
	go a.consumer.Start(a.workerCtx)

	a.initGRPCServer()
	a.initServer()

	return nil
}
