package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
	AppEnv      string `env:"APP_ENV" envDefault:"development"`
	JWTSecret   string `env:"JWT_SECRET"`
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

	if len(config.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 bytes")
	}

	return config, nil
}
