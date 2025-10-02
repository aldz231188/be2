package db

import (
	"be2/internal/domain"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewPGConfig,
		NewPool,
		fx.Annotate(
			NewRepo,
			fx.As(new(domain.AddressRepo)),
			// fx.As(new(domain.CustomerRepo)),
			// fx.As(new(domain.SupplierRepo)),
		),
	), fx.Invoke(),
)
