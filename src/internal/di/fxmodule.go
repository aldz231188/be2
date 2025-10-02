package di

import (
	"be2/internal/app"
	router "be2/internal/http"
	"be2/internal/http/v1/handlers"
	"be2/internal/infra/db"
	"go.uber.org/fx"
	"os"
)

var App = fx.Options(
	db.Module,
	fx.Provide(app.NewServiceImpl),
	fx.Provide(handlers.NewHandler),
	fx.Invoke(router.RegisterRoutes),
	fx.Invoke(func(g fx.DotGraph) {
		err := os.WriteFile("graph.dot", []byte(g), 0644)
		if err != nil {
			panic(err)
		}
	}),
)
