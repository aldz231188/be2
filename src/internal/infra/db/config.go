package db

import (
	"fmt"
	"os"
)

type Config struct {
	DSN string
}

// TODO: Не хардкодить DSN в коде: подними из env/файла в другом месте (например, infra/config).

func NewPGConfig() (Config, error) {
	cfg := Config{
		// DSN: "postgres://postgres:Qwaszx_1@localhost:5432/shopdb",
		DSN: os.Getenv("DATABASE_DSN"),
	}
	if cfg.DSN == "" {
		return cfg, fmt.Errorf("DATABASE_DSN is required")
	}
	return cfg, nil
}

// package config

// import (
// "fmt"
// "os"
// )

// type Config struct {
// Env string
// HTTPAddr string // ":8080"
// DSN string // "postgres://user:pass@localhost:5432/app?sslmode=disable"
// LogLevel string // "info|debug"
// }

// func FromEnv() (Config, error) {
// 	cfg := Config{
// 		// Env: getEnv("APP_ENV", "dev"),
// 		// HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
// 		DSN: os.Getenv("DATABASE_DSN"),
// 		// LogLevel: getEnv("LOG_LEVEL", "info"),
// 	}
// 	if cfg.DSN == "" {
// 		return cfg, fmt.Errorf("DATABASE_DSN is required")
// 	}
// 	return cfg, nil
// }

// func getEnv(k, def string) string {
// 	if v := os.Getenv(k); v != "" {
// 		return v
// 	}
// 	return def
// }
