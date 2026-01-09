package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	DB_URL string `env:"DB_URL"`
	Port   int    `env:"PORT" envDefault:"8000"`
}

func New() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("Failed to parse environment variables into config: %w", err)
	}

	return cfg, nil
}
