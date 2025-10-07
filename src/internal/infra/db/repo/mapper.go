package repo

import (
	"be2/internal/domain"
	store "be2/internal/infra/db/sqlc_generated"
)

func createAddressToRow(address domain.Address) store.CreateAddressParams {
	return store.CreateAddressParams{
		ID:      address.ID,
		Country: address.Country,
		City:    address.City,
		Street:  address.Street,
	}
}
func updateAddressToRow(address domain.Address) store.UpdateAddressParams {
	return store.UpdateAddressParams{
		ID:      address.ID,
		Country: address.Country,
		City:    address.City,
		Street:  address.Street,
	}
}

func createClientToRow(c domain.Client) store.CreateClientParams {
	return store.CreateClientParams{
		ID:         c.ID,
		ClientName: c.ClientName,
	}
}
