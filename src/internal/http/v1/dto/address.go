package dto

import (
	"be2/internal/domain"
)

type CreateAddressRequest struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Street  string `json:"street"`
}

type AddressResponse struct {
	ID      string `json:"id"`
	Country string `json:"name"`
	City    string `json:"surname"`
	Street  string `json:"gender"`
}

func (r CreateAddressRequest) ToDomainAdress() *domain.Adress {
	return &domain.Adress{
		Country: r.Country,
		City:    r.City,
		Street:  r.Street,
	}
}

func FromDomainAdress(c domain.Adress) AddressResponse {
	return AddressResponse{
		Country: c.Country,
		City:    c.City,
		Street:  c.Street,
	}
}
