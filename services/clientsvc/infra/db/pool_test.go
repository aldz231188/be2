package db

import (
	"be2/services/clientsvc/config"
	"net/url"
	"strings"
	"testing"
)

func TestBuildPostgresDSN(t *testing.T) {
	cfg := config.Config{
		Host:    "db.local",
		Port:    "5432",
		User:    "appuser",
		DB:      "service",
		SSLMode: "disable",
	}

	dsn := buildPostgresDSN(cfg, "s3cr3t")

	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}

	if u.Scheme != "postgres" {
		t.Fatalf("unexpected scheme: %s", u.Scheme)
	}
	if u.Host != "db.local:5432" {
		t.Fatalf("unexpected host: %s", u.Host)
	}
	if u.Path != "/service" {
		t.Fatalf("unexpected path: %s", u.Path)
	}
	if ssl := u.Query().Get("sslmode"); ssl != "disable" {
		t.Fatalf("unexpected sslmode: %s", ssl)
	}

	if username := u.User.Username(); username != "appuser" {
		t.Fatalf("unexpected username: %s", username)
	}
	if password, ok := u.User.Password(); !ok || password != "s3cr3t" {
		t.Fatalf("unexpected password: %s", password)
	}
}

func TestBuildPostgresDSN_EncodesCredentials(t *testing.T) {
	cfg := config.Config{
		Host:    "127.0.0.1",
		Port:    "5433",
		User:    "us er",
		DB:      "example",
		SSLMode: "verify-full",
	}

	password := "p@ss:w/or?d"
	dsn := buildPostgresDSN(cfg, password)

	if !strings.Contains(dsn, "us%20er") {
		t.Fatalf("username should be percent-encoded, got %s", dsn)
	}
	if !strings.Contains(dsn, "p%40ss%3Aw%2For%3Fd") {
		t.Fatalf("password should be percent-encoded, got %s", dsn)
	}

	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}

	if username := u.User.Username(); username != "us er" {
		t.Fatalf("decoded username mismatch: %s", username)
	}
	if decoded, ok := u.User.Password(); !ok || decoded != password {
		t.Fatalf("decoded password mismatch: %s", decoded)
	}
	if u.Path != "/example" {
		t.Fatalf("unexpected path: %s", u.Path)
	}
}
