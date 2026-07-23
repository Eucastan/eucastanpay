package bootstrap

import "github.com/Eucastan/eucastanpay/services/ledger/internal/repository/postgres"

func (a *App) initRepository() {
	a.repo = postgres.NewLedgerRepository(a.database.DB, a.telemetry)
}
