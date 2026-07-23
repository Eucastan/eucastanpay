package bootstrap

import "github.com/Eucastan/eucastanpay/services/admin/internal/repository/postgres"

func (a *App) initRepository() {
	a.repo = postgres.NewAdminRepository(a.database.DB)
}
