package dto

import (
	"be2/internal/app"
	"be2/internal/domain"
	"be2/internal/shared/date"
	"errors"

	// "time"
	"strings"
)

type CreateClientRequest struct {
	ClientName    string               `json:"client_name"`
	ClientSurname string               `json:"client_surname"`
	Birthday      date.DateOnly        `json:"birthday"`
	Gender        string               `json:"gender"`
	Address       CreateAddressRequest `json:"address"`
}

type ClientResponse struct {
	ID               string          `json:"id"`
	ClientName       string          `json:"client_name"`
	ClientSurname    string          `json:"client_surname"`
	Birthday         date.DateOnly   `json:"birthday"`
	Gender           string          `json:"gender"`
	RegistrationDate date.DateOnly   `json:"registration_date"`
	Address          AddressResponse `json:"address"`
}

func (r CreateClientRequest) ToDomainAddressClient() (app.CreateClientInput, error) {

	if gender, err := parseGender(r.Gender); err != nil {
		return app.CreateClientInput{}, err
	} else {

		return app.CreateClientInput{
			ClientName:    r.ClientName,
			ClientSurname: r.ClientSurname,
			Birthday:      r.Birthday,
			Gender:        gender,
			Address:       r.Address.ToDomainAddress(),
		}, nil
	}
}

// func FromDomainClient(c domain.Client) ClientResponse {
// 	return ClientResponse{
// 		ID:            c.ID.String(),
// 		ClientName:    c.ClientName,
// 		ClientSurname: c.ClientSurname,
// 		Address:       FromDomain(c.Address),
// 	}
// }

func parseGender(s string) (domain.Gender, error) {
	switch {
	case strings.EqualFold(s, "male"):
		return domain.MALE, nil
	case strings.EqualFold(s, "female"):
		return domain.FEMALE, nil
	default:
		return domain.MALE, errors.New("incorrect gender")

	}

}
