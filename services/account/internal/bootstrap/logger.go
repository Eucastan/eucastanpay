package bootstrap

import "github.com/Eucastan/eucastanpay/common/pkg/logger"

func (a *App) initLogger() {
	a.logger = logger.New(a.cfg.SharedCfg.LogLevel)
	a.logger.Info("Starting Account service...")
}
