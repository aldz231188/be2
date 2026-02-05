package config

import (
	"go.uber.org/fx"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DepTimeout      time.Duration // таймаут на один вызов зависимости
	ShutdownTimeout time.Duration
	AuthSvcAddr     string
	Host            string
	Port            string
	User            string
	DB              string
	SSLMode         string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		Host:            os.Getenv("DB_HOST"),
		Port:            os.Getenv("DB_PORT"),
		User:            os.Getenv("DB_USER"),
		DB:              os.Getenv("DB_NAME"),
		SSLMode:         os.Getenv("SSLMODE"),
		AuthSvcAddr:     os.Getenv("AUTH_SVC_ADDR_FULL"),
		DepTimeout:      msEnv("DEP_TIMEOUT_MS", 1500) * time.Millisecond,
		ShutdownTimeout: msEnv("SHUTDOWN_TIMEOUT_MS", 10_000) * time.Millisecond,
	}
	return cfg, nil
}

var Module = fx.Provide(LoadConfig, LoadSecrets)

func msEnv(k string, def int) time.Duration {
	v := os.Getenv(k)
	if v == "" {
		return time.Duration(def)
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return time.Duration(def)
	}
	return time.Duration(n)
}
