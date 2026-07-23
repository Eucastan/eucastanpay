package bootstrap

import "github.com/Eucastan/eucastanpay/services/ledger/internal/usecase/service"

func (a *App) initUseCase() {
	a.uc = service.NewLedgerUseCase(
		a.repo,
		a.telemetry,
		a.manager,
		a.logger,
	)
}
