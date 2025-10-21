package db

import (
	"context"
	"fmt"
	"time"

	"be2/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"net"
	"net/url"
)

func NewPool(lc fx.Lifecycle, cfg Config, sec *config.Secrets) (*pgxpool.Pool, error) {
	dsn := BuildPostgresDSN(cfg.Host, cfg.Port, cfg.User, sec.DBPassword, cfg.DB, cfg.SSLMode)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			c, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			if err := pool.Ping(c); err != nil {
				pool.Close()
				return fmt.Errorf("db ping failed: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func BuildPostgresDSN(host, port, user, pass, db, sslmode string) string {
	u := &url.URL{
		Scheme: "postgres",                   // или "postgresql"
		User:   url.UserPassword(user, pass), // ← корректное кодирование user:pass
		Host:   net.JoinHostPort(host, port),
		Path:   db, // "/app"
	}
	q := url.Values{}
	q.Set("sslmode", sslmode) // query-часть — уже по правилам query
	u.RawQuery = q.Encode()
	return u.String()
}
