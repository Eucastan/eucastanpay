package bootstrap

import "github.com/gin-gonic/gin"

func (a *App) initRouter() {
	a.router = gin.Default()
}
