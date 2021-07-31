package conf

import (
	"github.com/kelseyhightower/envconfig"
)

type BlizzApiCfg struct {
	RegionList   []string
	EuAPIUrl     string `envconfig:"BLIZZARD_EU_API_URL" default:"https://eu.api.blizzard.com"`
	UsAPIUrl     string `envconfig:"BLIZZARD_US_API_URL" default:"https://us.api.blizzard.com"`
	AUTHUrl      string `envconfig:"BLIZZARD_AUTH_URL" default:"https://us.battle.net/oauth/token"`
	ClientSecret string `envconfig:"BLIZZARD_CLIENT_SECRET" required:"true"`
	ClientID     string `envconfig:"BLIZZARD_CLIENT_ID" required:"true"`
	AuthTimeOut  int    `envconfig:"BLIZZARD_AUTH_TIMEOUT" default:"3"`
}

type Config struct {
	AppPort string `envconfig:"APP_PORT" required:"true"`
	LogLvl  string `envconfig:"LOG_LEVEL" default:"INFO"`
	BlizzApiCfg
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	cfg.RegionList = []string{"eu", "us"}

	if err := envconfig.Process("AUCTIONEER", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
