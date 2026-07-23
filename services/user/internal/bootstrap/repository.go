package bootstrap

import "github.com/Eucastan/eucastanpay/services/user/internal/repository/postgres"

func (a *App) initRepository() {
	a.userRepo = postgres.NewUserRepository(a.database.DB, a.telemetry)
	a.authRepo = postgres.NewAuthRepository(a.database.DB, a.telemetry)
	a.kycRepo = postgres.NewKYCRepository(a.database.DB, a.telemetry)
}
