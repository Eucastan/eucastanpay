package healthcheck

import (
	"context"
)

type Check interface {
	Name() string
	Check(ctx context.Context) Component
}

func (c *Checker) updateSummary(report *Report, component Component) {

	switch component.Status {

	case Healthy:
		report.Summary.Healthy++

	case Degraded:
		report.Summary.Degraded++

		if report.Status == Healthy {
			report.Status = Degraded
		}

	case Unhealthy:
		report.Summary.Unhealthy++
		report.Status = Unhealthy
	}
}

func (c *Checker) recoverComponent(check Check, report *Report) {

	if r := recover(); r != nil {
		component := Component{
			Name:   check.Name(),
			Status: Unhealthy,
			Error:  "panic during health check",
		}

		report.Components = append(
			report.Components,
			component,
		)

		c.updateSummary(
			report,
			component,
		)
	}
}
