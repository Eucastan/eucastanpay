package bootstrap

import "github.com/Eucastan/eucastanpay/services/transfer/internal/repository/postgres"

func (a *App) initRepository() {
	a.repo = postgres.NewTransferRepository(a.database.DB, a.telemetry)
}
