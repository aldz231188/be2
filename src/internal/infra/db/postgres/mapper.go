package db

import (
	"be2/internal/domain"
	store "be2/internal/infra/db/sqlc_generated"
)

func domaineToRow(address domain.Adress) store.AddAddressParams {
	return store.AddAddressParams{
		ID:      address.Id,
		Country: address.Country,
		City:    address.City,
		Street:  address.Street,
	}
}
