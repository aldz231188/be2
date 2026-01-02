package authclient

import (
	"os"
	"time"
)

type Config struct {
	Target  string
	Timeout time.Duration
}

func LoadConfig() Config {
	target := os.Getenv("AUTH_GRPC_TARGET")
	if target == "" {
		target = "auth:9090"
	}

	timeout := 3 * time.Second
	if v := os.Getenv("AUTH_GRPC_TIMEOUT"); v != "" {
		if parsed, err := time.ParseDuration(v); err == nil {
			timeout = parsed
		}
	}

	return Config{
		Target:  target,
		Timeout: timeout,
	}
}
