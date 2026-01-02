package config

import (
	"go.uber.org/fx"

	"os"
)

type Config struct {
	Host    string
	Port    string
	User    string
	DB      string
	SSLMode string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		Host:    os.Getenv("DB_HOST"),
		Port:    os.Getenv("DB_PORT"),
		User:    os.Getenv("DB_USER"),
		DB:      os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("SSLMODE"),
	}
	return cfg, nil
}

var Module = fx.Provide(LoadConfig, LoadSecrets)
