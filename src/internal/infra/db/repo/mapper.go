package repo

import (
	"be2/internal/domain"
	store "be2/internal/infra/db/sqlc_generated"
)

func domaineToRow(address domain.Address) store.CreateAddressParams {
	return store.CreateAddressParams{
		ID:      address.Id,
		Country: address.Country,
		City:    address.City,
		Street:  address.Street,
	}
}
