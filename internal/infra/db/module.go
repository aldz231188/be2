package db

import (
	"be2/internal/config"
	"be2/internal/domain"
	"be2/internal/infra/db/repo"
	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	fx.Provide(
		NewPool,
		fx.Annotate(
			repo.NewRepo,
			fx.As(new(domain.AddressRepo)),
			fx.As(new(domain.ClientRepo)),
			fx.As(new(domain.UserRepo)),
			fx.As(new(domain.SessionRepo)),
		// fx.As(new(domain.SupplierRepo)),
		),
	),
)
