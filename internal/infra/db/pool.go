package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"time"
)

func NewPool(lc fx.Lifecycle, dsn string) (*pgxpool.Pool, error) {

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
