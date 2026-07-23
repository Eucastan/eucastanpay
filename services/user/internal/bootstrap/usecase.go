package bootstrap

import (
	"github.com/Eucastan/eucastanpay/services/user/internal/usecase/service"
	"github.com/Eucastan/eucastanpay/services/user/internal/worker"
)

func (a *App) initUseCase() {
	publish := worker.NewPublishUserRegistration(a.userRepo)

	a.userUC = service.NewUserUseCase(
		a.userRepo,
		a.authRepo,
		a.telemetry,
		a.cfg,
		a.email,
		a.redis,
		publish,
	)

	a.kycUC = service.NewKYCUseCase(
		a.kycRepo,
		a.telemetry,
		a.cfg,
	)
}
