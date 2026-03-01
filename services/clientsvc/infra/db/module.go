package db

import (
	"be2/services/clientsvc/config"
	"be2/services/clientsvc/domain"
	"be2/services/clientsvc/infra/db/repo"
	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	fx.Provide(
		NewPool,
		fx.Annotate(
			repo.NewRepo,
			// fx.As(new(domain.AddressRepo)),
			fx.As(new(domain.ClientRepo)),
		),
	),
)
