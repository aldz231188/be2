package di

import (
	"be2/services/authsvc/internal/app"
	usecase "be2/services/authsvc/internal/app/usecase"
	client "be2/services/authsvc/internal/clients/user"
	router "be2/services/authsvc/internal/http"
	"be2/services/authsvc/internal/http/middleware"
	"be2/services/authsvc/internal/http/v1/handlers"
	"be2/services/authsvc/internal/infra/db"
	"log/slog"
	"os"

	"go.uber.org/fx"
)

var App = fx.Options(
	db.Module,
	fx.Provide(
		newLogger,
		usecase.NewClientUsecase,
		client.NewConn,
		client.NewService,
		fx.Annotate(
			app.NewAuthService,
			fx.As(new(app.AuthService)),
		),
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
