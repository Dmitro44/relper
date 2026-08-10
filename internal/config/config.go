package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
	AppEnv      string `env:"APP_ENV" envDefault:"development"`
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = env.Parse(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
