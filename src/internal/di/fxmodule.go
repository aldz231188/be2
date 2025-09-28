package di

import (
	as "be2/internal/api/http"
	"be2/internal/app"
	"be2/internal/infra/db"
	"go.uber.org/fx"
)

var App = fx.Options(
	db.Module,
	fx.Provide(app.NewClientServiceImpl),
	fx.Provide(as.NewHandler),
	fx.Invoke(as.RegisterRoutes),
)
