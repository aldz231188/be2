package di

import (
	"be2/internal/app"
	"be2/internal/grpc/authclient"
	"be2/internal/grpc/authserver"
	router "be2/internal/http"
	"be2/internal/http/middleware"
	"be2/internal/http/v1/handlers"
	"be2/internal/infra/db"
	"go.uber.org/fx"
	"log/slog"
	"os"
)

var BFFApp = fx.Options(
	db.Module,
	fx.Provide(
		newLogger,
		fx.Annotate(
			app.NewAddressService,
			fx.As(new(app.AddressService)),
		),
		fx.Annotate(
			app.NewClientService,
			fx.As(new(app.ClientService)),
		),
		fx.Annotate(
			authclient.NewAuthClient,
			fx.As(new(app.AuthService)),
		),
		authclient.LoadConfig,
		middleware.NewJWT,
	),
	fx.Provide(handlers.NewHandler),
	fx.Invoke(router.RegisterServer),
	fx.Provide(router.RegisterRoutes),
	fx.Invoke(func(g fx.DotGraph) {
		path := "/tmp/graph.dot"
		os.WriteFile(path, []byte(g), 0644)

	}),
)

var AuthApp = fx.Options(
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
