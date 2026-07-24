package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuditRoutes(api *gin.RouterGroup, deps Dependencies, h Handlers) {

	admin := AdminGroup(api, deps.Config)

	{
		admin.GET("/audits", h.Audit.GetAllAuditReads)
		admin.GET("/audits/search", h.Audit.SearchAuditLogs)
		admin.GET("/audits/:id", h.Audit.GetAuditRead)
	}
}
