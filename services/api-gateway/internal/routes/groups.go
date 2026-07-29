package routes

import (
	commonmw "github.com/Eucastan/eucastanpay/common/pkg/middleware"
	"github.com/Eucastan/eucastanpay/services/api-gateway/config"
	"github.com/gin-gonic/gin"
)

func UserGroup(api *gin.RouterGroup, cfg *config.Config) *gin.RouterGroup {

	g := api.Group("/")
	g.Use(commonmw.Auth(cfg.SharedCfg.JWTSecret))
	g.Use(commonmw.RequireRole("user"))

	return g
}

func AdminGroup(api *gin.RouterGroup, cfg *config.Config) *gin.RouterGroup {

	g := api.Group("/admin")
	g.Use(commonmw.AdminAuth(cfg.SharedCfg.JWTSecret))
	g.Use(commonmw.RequireAdminRole("admin", "super_admin"))

	return g
}

func SuperAdmin(api *gin.RouterGroup, cfg *config.Config) *gin.RouterGroup {

	g := api.Group("/admin")
	g.Use(commonmw.AdminAuth(cfg.SharedCfg.JWTSecret))
	g.Use(commonmw.RequireAdminRole("super_admin"))

	return g
}
