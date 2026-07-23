package bootstrap

import (
	"github.com/Eucastan/eucastanpay/common/pkg/logger"
)

func (a *App) initLogger() {
	log := logger.New(a.cfg.LogLevel)

	log.Info("Starting Admin service...")

	a.logger = log
}
