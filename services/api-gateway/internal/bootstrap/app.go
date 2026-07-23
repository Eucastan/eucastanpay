package bootstrap

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	manager "github.com/Eucastan/eucastanpay/common/pkg/grpc"
	"github.com/Eucastan/eucastanpay/common/pkg/grpc/discovery"
	"github.com/Eucastan/eucastanpay/common/pkg/healthcheck"
	"github.com/Eucastan/eucastanpay/common/pkg/logger"
	sharedmw "github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/Eucastan/eucastanpay/services/api-gateway/config"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/middleware"
	"github.com/Eucastan/eucastanpay/services/api-gateway/internal/ratelimiter"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type App struct {
	cfg          *config.Config
	router       *gin.Engine
	server       *http.Server
	logger       *logrus.Logger
	redis        *redis.Client
	telemetry    *telemetry.Telemetry
	rateLimiter  *ratelimiter.RedisLimiter
	manager      *manager.Manager
	health       *healthcheck.Checker
	gateways     *Gateways
	applications *Applications
	handlers     *Handlers
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
	go func() {
		a.logger.Infof(
			"Gateway listening on %s",
			a.cfg.HTTPPort,
		)

		if err := a.server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			a.logger.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	a.logger.Info("shutdown signal received")

	return a.shutdown(context.Background())
}

func (a *App) bootstrap() error {
	a.initLogger()

	if err := a.initTelemetry(); err != nil {
		return err
	}

	registry := discovery.NewStaticRegistry(
		map[string]string{
			"user":     a.cfg.UserGRPCAddr,
			"account":  a.cfg.AccountGRPCAddr,
			"transfer": a.cfg.TransferGRPCAddr,
			"ledger":   a.cfg.LedgerGRPCAddr,
			"audit":    a.cfg.AuditGRPCAddr,
		},
	)

	a.manager = manager.NewManager(registry)

	if err := a.initClients(); err != nil {
		return err
	}

	a.initGateways()
	a.initApplications()
	a.initHandlers()
	a.initRouter()
	a.registerMiddleware()
	a.initHealth()
	a.registerRoutes()
	a.initServer()

	return nil
}

func (a *App) initLogger() {
	a.logger = logger.New(
		a.cfg.SharedCfg.LogLevel,
	)
}

func (a *App) initTelemetry() error {
	tracer := otel.Tracer("gateway")
	meter := otel.Meter("gateway")

	tm, err := telemetry.New(tracer, meter, a.logger)
	if err != nil {
		return err
	}

	a.telemetry = tm

	return nil
}

func (a *App) registerMiddleware() {

	a.router.Use(

		middleware.Recovery(),
		middleware.RequestID(),
		sharedmw.CorrelationMiddleware(),
		middleware.Logger(a.logger),
		middleware.CORS(),
	)
}

func (a *App) initServer() {

	a.server = &http.Server{

		Addr: ":" + a.cfg.HTTPPort,

		Handler: a.router,
	}
}
