package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadSecretFileOrEnv_PrefersFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secret.txt")
	if err := os.WriteFile(path, []byte("from-file"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	t.Setenv("DB_PASSWORD_FILE", path)
	t.Setenv("DB_PASSWORD", "from-env")

	secret, err := readSecretFileOrEnv("DB_PASSWORD_FILE", "DB_PASSWORD")
	if err != nil {
		t.Fatalf("readSecretFileOrEnv returned error: %v", err)
	}
	if secret != "from-file" {
		t.Fatalf("expected file content, got %q", secret)
	}
}

func TestReadSecretFileOrEnv_EnvFallback(t *testing.T) {
	t.Setenv("DB_PASSWORD_FILE", "")
	t.Setenv("DB_PASSWORD", "from-env")

	secret, err := readSecretFileOrEnv("DB_PASSWORD_FILE", "DB_PASSWORD")
	if err != nil {
		t.Fatalf("readSecretFileOrEnv returned error: %v", err)
	}
	if secret != "from-env" {
		t.Fatalf("expected env content, got %q", secret)
	}
}

func TestReadSecretFileOrEnv_Missing(t *testing.T) {
	t.Setenv("DB_PASSWORD_FILE", "")
	t.Setenv("DB_PASSWORD", "")

	if _, err := readSecretFileOrEnv("DB_PASSWORD_FILE", "DB_PASSWORD"); err == nil {
		t.Fatal("expected error when neither file nor env is set")
	}
}

func TestLoadSecrets_FromFiles(t *testing.T) {
	dir := t.TempDir()

	dbPath := filepath.Join(dir, "dbpass.txt")
	jwtPath := filepath.Join(dir, "jwt.txt")

	if err := os.WriteFile(dbPath, []byte("db-from-file"), 0o600); err != nil {
		t.Fatalf("write db file: %v", err)
	}
	if err := os.WriteFile(jwtPath, []byte("jwt-from-file"), 0o600); err != nil {
		t.Fatalf("write jwt file: %v", err)
	}

	t.Setenv("DB_PASSWORD_FILE", dbPath)
	t.Setenv("JWT_PRIVATE_KEY_FILE", jwtPath)
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("JWT_PRIVATE_KEY", "")

	secrets, err := LoadSecrets()
	if err != nil {
		t.Fatalf("LoadSecrets returned error: %v", err)
	}

	if secrets.DBPassword != "db-from-file" || secrets.JWTPrivateKey != "jwt-from-file" {
		t.Fatalf("unexpected secrets: %+v", secrets)
	}
}

func TestLoadSecrets_FromEnv(t *testing.T) {
	t.Setenv("DB_PASSWORD_FILE", "")
	t.Setenv("JWT_PRIVATE_KEY_FILE", "")
	t.Setenv("DB_PASSWORD", "db-env")
	t.Setenv("JWT_PRIVATE_KEY", "jwt-env")

	secrets, err := LoadSecrets()
	if err != nil {
		t.Fatalf("LoadSecrets returned error: %v", err)
	}

	if secrets.DBPassword != "db-env" || secrets.JWTPrivateKey != "jwt-env" {
		t.Fatalf("unexpected secrets: %+v", secrets)
	}
}

func TestLoadSecrets_MissingValues(t *testing.T) {
	t.Setenv("DB_PASSWORD_FILE", "")
	t.Setenv("JWT_PRIVATE_KEY_FILE", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("JWT_PRIVATE_KEY", "")

	if _, err := LoadSecrets(); err == nil {
		t.Fatal("expected error when required secrets are missing")
	}
}
