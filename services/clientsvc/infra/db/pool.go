package db

import (
	"be2/services/clientsvc/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"net"
	"net/url"
	"time"
)

func NewPool(lc fx.Lifecycle, cfg config.Config, sec *config.Secrets) (*pgxpool.Pool, error) {
	dsn := buildPostgresDSN(cfg, sec.DBPassword)
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

func buildPostgresDSN(cfg config.Config, pass string) string { //ввынести?
	u := &url.URL{
		Scheme: "postgres",                       // или "postgresql"
		User:   url.UserPassword(cfg.User, pass), // ← корректное кодирование user:pass
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   cfg.DB,
	}
	q := url.Values{}
	q.Set("sslmode", cfg.SSLMode) // query-часть — уже по правилам query
	u.RawQuery = q.Encode()
	return u.String()
}
