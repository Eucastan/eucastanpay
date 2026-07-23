package bootstrap

import "github.com/Eucastan/eucastanpay/services/account/internal/repository/postgres"

func (a *App) initRepository() {
	a.repo = postgres.NewAccountRepository(a.database.DB, a.telemetry, a.logger)
}
