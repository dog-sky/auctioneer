package conf

import (
	"github.com/kelseyhightower/envconfig"
	"net/url"
)

type BlizzApiCfg struct {
	APIUrl       *url.URL `envconfig:"BLIZZARD_API_URL" required:"true"`
	AUTHUrl      *url.URL `envconfig:"BLIZZARD_AUTH_URL" required:"true"`
	ClientSecret string   `envconfig:"BLIZZARD_CLIENT_SECRET" required:"true"`
	ClientID     string   `envconfig:"BLIZZARD_CLIENT_ID" required:"true"`
}

type Config struct {
	AppPort string `envconfig:"APP_PORT" required:"true"`
	LogLvl  string `envconfig:"LOG_LEVEL" default:"INFO"`
	BlizzApiCfg
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	if err := envconfig.Process("AUCTIONEER", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
