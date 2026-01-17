package config

import (
	// "context"
	"errors"
	"os"
	// "go.uber.org/fx"
)

type Secrets struct {
	DBPassword string
}

func readSecretFileOrEnv(fileKey, envKey string) (string, error) {
	if p := os.Getenv(fileKey); p != "" {
		b, err := os.ReadFile(p)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	if v := os.Getenv(envKey); v != "" {
		return v, nil
	}
	return "", errors.New("missing secret: " + fileKey + " or " + envKey)
}

func LoadSecrets() (*Secrets, error) {
	dbPass, err := readSecretFileOrEnv("DB_PASSWORD_FILE", "DB_PASSWORD")
	if err != nil {
		return nil, err
	}
	return &Secrets{
		DBPassword: dbPass,
	}, nil
}
