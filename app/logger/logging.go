package logging

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(logLvl string) (*Logger, error) {
	lvl, err := logrus.ParseLevel(logLvl)
	if err != nil {
		return nil, err
	}

	baseLogger := logrus.New()
	logger := &Logger{baseLogger}
	logger.SetLevel(lvl)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger, nil
}
