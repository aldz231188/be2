package db

import (
	"be2/internal/domain"
	sqlc "be2/internal/infra/db/sqlc_generated"
)

func domaineToRow(address *domain.Adress) *sqlc.AddAddressParams {
	return &sqlc.AddAddressParams{
		ID:      address.Id,
		Country: address.Country,
		City:    address.City,
		Street:  address.Street,
	}
}
