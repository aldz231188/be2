package auth

import (
	"context"
	// "time"

	"be2/internal/config"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Conn struct{ *grpc.ClientConn }

func NewConn(lc fx.Lifecycle, cfg config.Config) (Conn, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	cc, err := grpc.NewClient(cfg.AuthSvcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // dev-only
		// grpc.WithBlock(),
	)
	if err != nil {
		return Conn{}, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error { return cc.Close() },
	})
	return Conn{ClientConn: cc}, nil
}
