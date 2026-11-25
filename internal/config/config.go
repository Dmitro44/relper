package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	AppEnv      string `env:"APP_ENV" envDefault:"development"`
	JWTSecret   string `env:"JWT_SECRET,required"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	config := new(Config)
	if err := env.Parse(config); err != nil {
		return nil, err
	}

	if len(config.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 bytes")
	}

	return config, nil
}
