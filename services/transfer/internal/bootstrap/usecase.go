package bootstrap

import "github.com/Eucastan/eucastanpay/services/transfer/internal/usecase/service"

func (a *App) initUseCase() {
	a.uc = service.NewTransferUseCase(
		a.repo,
		a.manager,
		a.publish,
		a.telemetry,
		a.logger,
	)
}
