package bootstrap

import "github.com/Eucastan/eucastanpay/services/audit/internal/usecase/service"

func (a *App) initUseCase() {
	a.uc = service.NewAuditUseCase(a.repo, a.telemetry)
}
