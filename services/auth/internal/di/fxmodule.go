package di

import (
	"be2/services/auth/internal/app"
	"be2/services/auth/internal/grpc/authserver"
	"be2/services/auth/internal/infra/db"
	"go.uber.org/fx"
	"log/slog"
	"os"
)

var App = fx.Options(
	db.Module,
	fx.Provide(
		newLogger,
		fx.Annotate(
			app.NewAuthService,
			fx.As(new(app.AuthService)),
		),
		authserver.LoadConfig,
		authserver.NewServer,
	),
	fx.Invoke(authserver.Register),
)

func newLogger() *slog.Logger { // вынести
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	return slog.New(handler)
}
