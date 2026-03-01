package di

import (
	usecase "be2/internal/app/usecase"
	auth "be2/internal/clients/auth"
	client "be2/internal/clients/client"
	router "be2/internal/http"
	"be2/internal/http/v1/handlers"

	// "example.com/bff/internal/clients/auth"
	"be2/internal/config"
	"be2/internal/grpcutil"
	// "example.com/bff/internal/http/handlers"
	"be2/internal/http/middleware"
	// "example.com/bff/internal/http/server"
	"log/slog"
	"os"

	"go.uber.org/fx"
)

var App = fx.Options(
	grpcutil.Module,
	middleware.Module,
	config.Module,
	fx.Provide(
		newLogger,
		usecase.NewClientUsecase,
		usecase.NewAuthUsecase,
		client.NewConn,
		client.NewService,
		auth.NewConn,
		auth.NewService,
		// fx.Annotate(
		// 	app.NewAuthService,
		// 	fx.As(new(app.AuthService)),
		// ),
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
