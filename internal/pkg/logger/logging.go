package logging

import (
	"github.com/sirupsen/logrus"
)

func NewLogger(logLvl string) (*logrus.Logger, error) {
	lvl, err := logrus.ParseLevel(logLvl)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(lvl)
	logger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})

	return logger, nil
}
