package bootstrap

import "github.com/Eucastan/eucastanpay/services/account/internal/usecase/service"

func (a *App) initUseCase() {
	a.uc = service.NewAccountUseCase(
		a.repo,
		a.telemetry,
		a.logger,
	)
}
