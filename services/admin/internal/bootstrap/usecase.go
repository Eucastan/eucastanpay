package bootstrap

import "github.com/Eucastan/eucastanpay/services/admin/internal/usecase/service"

func (a *App) initUseCase() {
	a.uc = service.NewAdminUseCase(a.repo, a.cfg, a.logger)
}
