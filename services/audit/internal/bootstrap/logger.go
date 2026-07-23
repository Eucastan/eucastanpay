package bootstrap

import "github.com/Eucastan/eucastanpay/common/pkg/logger"

func (a *App) initLogger() {
	a.logger = logger.New(a.cfg.LogLevel)
	a.logger.Info("Starting Audit service...")
}
