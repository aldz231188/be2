package dto

import (
	"be2/internal/domain"
	// "be2/internal/shared/date"
	// "time"
)

type CreateClientRequest struct {
	ClientName    string `json:"client_name"`
	ClientSurname string `json:"client_surname"`
	// Birthday      date.DateOnly        `json:"birthday"`
	// Gender        string               `json:"gender"`
	Address CreateAddressRequest `json:"address"`
}

type ClientResponse struct {
	ID            string `json:"id"`
	ClientName    string `json:"client_name"`
	ClientSurname string `json:"client_surname"`
	// Birthday         date.DateOnly   `json:"birthday"`
	// Gender           string          `json:"gender"`
	// RegistrationDate date.DateOnly   `json:"registration_date"`
	Address AddressResponse `json:"address"`
}

func (r CreateClientRequest) ToDomainClient() *domain.Client {
	return &domain.Client{
		ClientName:    r.ClientName,
		ClientSurname: r.ClientSurname,
		// Birthday:      r.Birthday,
		// Gender:        r.Gender,
		Address: r.ToDomainClient().Address,
	}
}

func FromDomainClient(c domain.Client) ClientResponse {
	return ClientResponse{
		ID:            c.Id.String(),
		ClientName:    c.ClientName,
		ClientSurname: c.ClientSurname,
		Address:       FromDomain(c.Address),
	}
}
