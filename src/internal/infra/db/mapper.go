package db

import (
	"be2/internal/domain"
)

func domaineToRow(address *domain.Adress) *AddAddressParams {
	return &AddAddressParams{
		ID: address.Id,
	}
}
