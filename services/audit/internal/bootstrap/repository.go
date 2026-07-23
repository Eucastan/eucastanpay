package bootstrap

import (
	"github.com/Eucastan/eucastanpay/services/audit/internal/repository/postgres"
)

func (a *App) initRepository() {
	a.repo = postgres.NewAuditRepository(a.database.DB, a.telemetry)
}
