package dto

import (
	"be2/internal/app"
	"be2/internal/domain"
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

func (r CreateAddressRequest) ToDomain() app.CreateInput {
	return app.CreateInput{
		Country: r.Country,
		City:    r.City,
		Street:  r.Street,
	}
}

func FromDomain(c domain.Address) AddressResponse {
	return AddressResponse{
		Country: c.Country,
		City:    c.City,
		Street:  c.Street,
	}
}
