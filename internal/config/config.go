package config

import (
	"go.uber.org/fx"

	"os"
)

type Config struct {
	Host    string
	Port    string
	User    string
	Pass    string
	DB      string
	SSLMode string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		// DSN: "postgres://postgres:Qwaszx_1@localhost:5432/shopdb",
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		User: os.Getenv("DB_USER"),
		// Pass:    pass.DBPassword,
		DB:      os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("SSLMODE"),
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

var Module = fx.Provide(LoadConfig, LoadSecrets)
