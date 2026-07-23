package healthcheck

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (c *Checker) report(ctx context.Context, checks []Check) Report {
	report := Report{
		Service:    c.service,
		Version:    c.version,
		Status:     Healthy,
		StartedAt:  c.started,
		Uptime:     time.Since(c.started).String(),
		Timestamp:  time.Now().UTC(),
		Components: []Component{},
	}

	for _, check := range checks {
		func(ch Check) {
			defer c.recoverComponent(ch, &report)
			component := ch.Check(ctx)
			report.Components = append(report.Components, component)
			c.updateSummary(&report, component)

		}(check)
	}

	return report
}

func (c *Checker) Report(ctx context.Context) Report {
	return c.report(ctx, c.checks)
}

func (c *Checker) ReadyReport(ctx context.Context) Report {
	return c.report(ctx, c.readiness)
}

func (c *Checker) LiveReport(ctx context.Context) Report {
	return c.report(ctx, c.liveness)
}

func (c *Checker) respond(g *gin.Context, report Report) {

	status := http.StatusOK

	if report.Status == Unhealthy {
		status = http.StatusServiceUnavailable
	}

	g.JSON(status, report)
}

func (c *Checker) Health(g *gin.Context) {
	report := c.Report(g.Request.Context())
	c.respond(g, report)
}

func (c *Checker) Readiness(g *gin.Context) {
	if !c.IsReady() {

		g.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
		})

		return
	}

	report := c.ReadyReport(g.Request.Context())
	c.respond(g, report)
}

func (c *Checker) Liveness(g *gin.Context) {
	report := c.LiveReport(g.Request.Context())
	report.Status = "alive"

	c.respond(g, report)
}
