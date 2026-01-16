package di

import (
	"be2/services/clientsvc/app"
	"be2/services/clientsvc/infra/db"
	grpc "be2/services/clientsvc/transport"
	"be2/services/clientsvc/transport/grpc/handlers"
	"go.uber.org/fx"
	"log/slog"
	"os"
)

var App = fx.Options(
	db.Module,
	fx.Provide(
		newLogger,
		// fx.Annotate(
		// 	app.NewAddressService,
		// 	fx.As(new(app.AddressService)),
		// ),
		fx.Annotate(
			app.NewClientService,
			fx.As(new(app.ClientService)),
		),
		// fx.Annotate(
		// 	app.NewAuthService,
		// 	fx.As(new(app.AuthService)),
		// ),
	),
	fx.Provide(handlers.NewHandler),
	fx.Provide(grpc.NewGRPCServer,
		grpc.NewListener),
	fx.Invoke(grpc.Start),

	// fx.Provide(router.RegisterRoutes),
	// fx.Invoke(func(g fx.DotGraph) {
	// 	path := "/tmp/graph.dot"
	// 	os.WriteFile(path, []byte(g), 0644)

	// }),
)

func newLogger() *slog.Logger { // вынести
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	return slog.New(handler)
}
