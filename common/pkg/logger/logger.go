package logger

import "github.com/sirupsen/logrus"

func New(level string) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil
	}

	log.SetLevel(lvl)

	return log
}
