package authserver

import (
	"os"
)

type Config struct {
	Addr string
}

func LoadConfig() Config {
	addr := os.Getenv("AUTH_GRPC_ADDR")
	if addr == "" {
		addr = ":9090"
	}

	return Config{Addr: addr}
}
