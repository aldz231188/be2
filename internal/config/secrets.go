package config

import (
	// "context"
	"errors"
	"os"
	// "go.uber.org/fx"
)

type Secrets struct {
	DBPassword string
	// JWTPrivateKeyPEM  []byte
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

	// jwtPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	// if jwtPath == "" {
	// 	return nil, errors.New("missing secret: JWT_PRIVATE_KEY_PATH")
	// }
	// key, err := os.ReadFile(jwtPath)
	// if err != nil { return nil, err }

	return &Secrets{
		DBPassword: dbPass,
		// JWTPrivateKeyPEM: key,
	}, nil
}
