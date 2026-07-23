package bootstrap

import (
	"time"

	"github.com/Eucastan/eucastanpay/services/admin/internal/worker"
)

func (a *App) initOutboxWorker() {
	a.worker = worker.NewOutboxWorker(
		a.database.DB,
		a.publish,
		a.logger,
		2*time.Second,
	)
}
