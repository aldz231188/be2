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
	ClientSvcAddr   string
	AuthSvcAddr     string

	HTTPAddr string // ":8080"
	// AuthGRPCAddr  string // "auth:50051" или "127.0.0.1:50051"
	JWKSURL        string // "https://auth.example.com/.well-known/jwks.json"
	JWTIssuer      string // "auth.example.com"
	JWTAudience    string // опционально
	CookieName     string // "refresh_token"
	CookiePath     string // "/auth/refresh"
	CookieSecure   bool   // true на проде, false на локалке без TLS
	CookieSameSite string // "Lax"|"Strict"|"None"
}

func LoadConfig() (Config, error) {
	cfg := Config{
		HTTPAddr: os.Getenv("HTTP_PORT"),

		ClientSvcAddr:   os.Getenv("CLIENT_SVC_ADDR_FULL"),
		AuthSvcAddr:     os.Getenv("AUTH_SVC_ADDR_FULL"),
		DepTimeout:      msEnv("DEP_TIMEOUT_MS", 1500) * time.Millisecond,
		ShutdownTimeout: msEnv("SHUTDOWN_TIMEOUT_MS", 10_000) * time.Millisecond,

		JWKSURL:        os.Getenv("JWKS_URL"),
		JWTIssuer:      os.Getenv("JWT_ISS"),
		JWTAudience:    os.Getenv("JWT_AUD"),
		CookieName:     os.Getenv("COOKIE_NAME"),
		CookiePath:     os.Getenv("COOKIE_PATH"),
		CookieSecure:   getEnvBool("COOKIE_SECURE", false),
		CookieSameSite: os.Getenv("COOKIE_SAMESITE"),
	}
	return cfg, nil
}

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

func getEnvBool(k string, def bool) bool {
	if v := os.Getenv(k); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}

var Module = fx.Provide(LoadConfig, LoadSecrets)
