package bootstrap

import "github.com/Eucastan/eucastanpay/services/user/internal/usecase"

func (a *App) initEmail() {
	a.email = usecase.NewEmailService(a.cfg)
}
