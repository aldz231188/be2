package di

import (
	"be2/services/bff/internal/app"
	"be2/services/bff/internal/grpc/authclient"
	router "be2/services/bff/internal/http"
	"be2/services/bff/internal/http/middleware"
	"be2/services/bff/internal/http/v1/handlers"
	"be2/services/bff/internal/infra/db"
	"go.uber.org/fx"
	"log/slog"
	"os"
)

var App = fx.Options(
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

func newLogger() *slog.Logger { // вынести
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	return slog.New(handler)
}
