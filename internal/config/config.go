package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() (*Config, error) {
	config := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}

	if config.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL environment variable missing or empty")
	}
	if config.Port == "" {
		return nil, errors.New("PORT environment variable missing or empty")
	}

	return config, nil
}
