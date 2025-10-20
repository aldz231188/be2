package db

import (
	"be2/internal/domain"
	"be2/internal/infra/db/repo"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewPGConfig,
		NewPool,
		fx.Annotate(
			repo.NewRepo,
			fx.As(new(domain.AddressRepo)),
			fx.As(new(domain.ClientRepo)),
			// fx.As(new(domain.SupplierRepo)),
		),
	),
)
