package conf

import (
	"github.com/caarlos0/env/v6"
)

type App struct {
	PrometheusBind      string `env:"PROMETHEUS_BIND" envDefault:":2112"`
	HttpAddr            string `env:"HTTP_ADDR" envDefault:":9090"`
	PostgresDSN         string `env:"PG_DSN"`
	CustomOAuthRedirect string `env:"CUSTOM_OAUTH_REDIRECT"`
}

func ParseEnv() (*App, error) {
	cfg := App{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
