package dto

import (
	"be2/internal/app"
	"be2/internal/domain"
	"strings"
)

type CreateAddressRequest struct {
	Country string `json:"country" validate:"required"`
	City    string `json:"city"`
	Street  string `json:"street"`
}

type AddressResponse struct {
	ID      string `json:"id"`
	Country string `json:"name"`
	City    string `json:"surname"`
	Street  string `json:"gender"`
}

func (r CreateAddressRequest) ToDomainAddress() (app.CreateAddressInput, error) {
	errs := domain.NewValidationErrors()

	country := strings.TrimSpace(r.Country)
	if country == "" {
		errs.Add("address.country", "is required")
	}

	city := strings.TrimSpace(r.City)
	if city == "" {
		errs.Add("address.city", "is required")
	}

	street := strings.TrimSpace(r.Street)
	if street == "" {
		errs.Add("address.street", "is required")
	}

	if errs.HasErrors() {
		return app.CreateAddressInput{}, errs
	}

	return app.CreateAddressInput{
		Country: country,
		City:    city,
		Street:  street,
	}, nil
}

func FromDomainAddress(c domain.Address) AddressResponse {
	return AddressResponse{
		Country: c.Country,
		City:    c.City,
		Street:  c.Street,
	}
}
