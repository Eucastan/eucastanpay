package bootstrap

import (
	"time"

	"github.com/Eucastan/eucastanpay/services/ledger/internal/worker"
)

func (a *App) initOutboxWorker() {
	w := worker.NewOutboxWorker(
		a.database.DB,
		a.publish,
		a.logger,
		2*time.Second,
	)

	a.worker = w
}
