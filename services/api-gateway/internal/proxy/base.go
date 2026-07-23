package proxy

import (
	"github.com/Eucastan/eucastanpay/common/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

type Base struct {
	Logger    *logrus.Logger
	Telemetry *telemetry.Telemetry
}

func NewBase(
	logger *logrus.Logger,
	telemetry *telemetry.Telemetry,
) *Base {
	return &Base{
		Logger:    logger,
		Telemetry: telemetry,
	}
}
