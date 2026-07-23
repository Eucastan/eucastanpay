package bootstrap

import "github.com/Eucastan/eucastanpay/services/audit/internal/infra/database"

func (a *App) initDatabase() {
	a.database = database.NewPostgresDB(a.cfg, a.logger)
}
