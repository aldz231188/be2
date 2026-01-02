package dto

import (
	"be2/services/bff/internal/app"
	"be2/services/bff/internal/domain"
	"be2/services/bff/internal/shared/date"
	"strings"
	"time"
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
	errs := domain.NewValidationErrors()

	name := strings.TrimSpace(r.ClientName)
	if name == "" {
		errs.Add("client_name", "is required")
	}

	surname := strings.TrimSpace(r.ClientSurname)
	if surname == "" {
		errs.Add("client_surname", "is required")
	}

	if r.Birthday.Time.IsZero() {
		errs.Add("birthday", "is required")
	} else if r.Birthday.Time.After(time.Now().UTC()) {
		errs.Add("birthday", "cannot be in the future")
	}

	gender, err := parseGender(r.Gender)
	if err != nil {
		if validationErrs, ok := err.(*domain.ValidationErrors); ok {
			errs.Merge(validationErrs)
		} else {
			return app.CreateClientInput{}, err
		}
	}

	address, err := r.Address.ToDomainAddress()
	if err != nil {
		if validationErrs, ok := err.(*domain.ValidationErrors); ok {
			errs.Merge(validationErrs)
		} else {
			return app.CreateClientInput{}, err
		}
	}

	if errs.HasErrors() {
		return app.CreateClientInput{}, errs
	}

	return app.CreateClientInput{
		ClientName:    name,
		ClientSurname: surname,
		Birthday:      r.Birthday,
		Gender:        gender,
		Address:       address,
	}, nil
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
	value := strings.TrimSpace(s)
	switch {
	case strings.EqualFold(value, "male"):
		return domain.MALE, nil
	case strings.EqualFold(value, "female"):
		return domain.FEMALE, nil
	case value == "":
		errs := domain.NewValidationErrors()
		errs.Add("gender", "is required")
		return domain.MALE, errs
	default:
		errs := domain.NewValidationErrors()
		errs.Add("gender", "must be either male or female")
		return domain.MALE, errs
	}
}
