package repo

import (
	"be2/services/clientsvc/domain"
	store "be2/services/clientsvc/infra/db/sqlc_generated"
)

// func createAddressToRow(address domain.Address) store.CreateAddressParams {
// 	return store.CreateAddressParams{
// 		ID:      address.ID,
// 		Country: address.Country,
// 		City:    address.City,
// 		Street:  address.Street,
// 	}
// }
// func updateAddressToRow(address domain.Address) store.UpdateAddressParams {
// 	return store.UpdateAddressParams{
// 		ID:      address.ID,
// 		Country: address.Country,
// 		City:    address.City,
// 		Street:  address.Street,
// 	}
// }

func createClientToRow(c domain.Client) store.CreateClientParams {
	var gender string
	// switch c.Gender {
	// case domain.FEMALE:
	// 	gender = "female"
	// case domain.MALE:
	// 	gender = "male"
	// }
	return store.CreateClientParams{
		ID:            c.ID,
		ClientName:    c.ClientName,
		ClientSurname: c.ClientSurname,
		// Birthday:      c.Birthday,
		Gender: gender,
		// AddressID:     c.Address,
	}
}
