package healthcheck

import (
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type Checker struct {
	service   string
	version   string
	started   time.Time
	logger    *logrus.Logger
	ready     atomic.Bool
	checks    []Check
	readiness []Check
	liveness  []Check
}

func New(service, version string, logger *logrus.Logger) *Checker {
	checker := &Checker{
		service:   service,
		version:   version,
		started:   time.Now(),
		logger:    logger,
		checks:    []Check{},
		readiness: []Check{},
		liveness:  []Check{},
	}

	checker.ready.Store(true)
	return checker
}

func (c *Checker) Add(checks ...Check) {

	for _, check := range checks {
		if check == nil {
			c.logger.Warn("nil health check ignored")
			continue
		}

		c.checks = append(c.checks, check)

		c.logger.Infof(
			"registered health check: %s",
			check.Name(),
		)
	}
}

func (c *Checker) AddReadiness(checks ...Check) {
	c.readiness = append(c.readiness, checks...)
}

func (c *Checker) AddLiveness(checks ...Check) {
	c.liveness = append(c.liveness, checks...)
}

func (c *Checker) SetReady(v bool) {
	c.ready.Store(v)
}

func (c *Checker) IsReady() bool {
	return c.ready.Load()
}
