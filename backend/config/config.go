package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	DB_URL string `env:"DB_URL"`
	Port   int    `env:"PORT" envDefault:"8000"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
