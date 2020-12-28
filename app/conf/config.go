package conf

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppPort string `envconfig:"APP_PORT"`
	LogLvl  string `envconfig:"LOG_LEVEL" default:"INFO"`
	EnvLvl  string `envconfig:"ENV" default:"dev" desc:"Описание среды окружения"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	if err := envconfig.Process("AUCTIONEER", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
