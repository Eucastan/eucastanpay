package bootstrap

import (
	"net/http"
	"time"
)

func (a *App) initServer() {
	a.server = &http.Server{
		Addr:         a.cfg.HTTPPort,
		Handler:      a.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
