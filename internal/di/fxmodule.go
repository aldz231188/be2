package di

import (
	"be2/internal/app"
	router "be2/internal/http"
	"be2/internal/http/middleware"
	"be2/internal/http/v1/handlers"
	"be2/internal/infra/db"
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
			app.NewAuthService,
			fx.As(new(app.AuthService)),
		),
		middleware.NewJWT,
	),
	fx.Provide(handlers.NewHandler),
	fx.Invoke(router.NewServer),
	fx.Provide(router.RegisterRoutes),
	// fx.Invoke(func(g fx.DotGraph) {
	// 	err := os.WriteFile("graph.dot", []byte(g), 0644)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }),
)

func newLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	return slog.New(handler)
}
