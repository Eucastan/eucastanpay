package bootstrap

import (
	"context"
	"time"

	adminkafka "github.com/Eucastan/eucastanpay/common/pkg/kafka/admin"
)

func (a *App) initTopics() error {
	a.topicsInitializer = adminkafka.NewInitializer(a.cfg.SharedCfg.Kafka)

	a.logger.Info("Initializing Kafka topics...")

	var err error

	for i := 0; i < 10; i++ {
		err = a.topicsInitializer.EnsureTopics(context.Background())
		if err == nil {
			a.logger.Info("Kafka topics are ready.")
			return nil
		}

		a.logger.Warnf(
			"Kafka not ready (attempt %d/10): %v",
			i+1,
			err,
		)

		time.Sleep(5 * time.Second)
	}

	return err
}
