package startup

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Task interface {
	Name() string
	Run(ctx context.Context) error
}

type Validator struct {
	logger *logrus.Logger
	tasks  []Task
}

func New(logger *logrus.Logger) *Validator {
	return &Validator{
		logger: logger,
	}
}

func (v *Validator) Add(tasks ...Task) {
	for _, task := range tasks {
		if task == nil {
			continue
		}

		v.tasks = append(v.tasks, task)
	}
}

func (v *Validator) Validate(ctx context.Context) error {

	v.logger.Info("running startup validation...")

	for _, task := range v.tasks {

		start := time.Now()

		v.logger.Infof(
			"checking %s...",
			task.Name(),
		)

		if err := task.Run(ctx); err != nil {

			v.logger.Errorf(
				"%s failed (%s): %v",
				task.Name(),
				time.Since(start),
				err,
			)

			return fmt.Errorf(
				"%s startup failed: %w",
				task.Name(),
				err,
			)
		}

		v.logger.Infof(
			"%s healthy (%s)",
			task.Name(),
			time.Since(start),
		)
	}

	v.logger.Info("startup validation completed")

	return nil
}
