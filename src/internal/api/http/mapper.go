package http

import (
	"be2/internal/domain"
)

func ToDomain(rowAddress Address) *domain.Adress {
	return &domain.Adress{
		Country: rowAddress.Country,
		City:    rowAddress.City,
		Street:  rowAddress.Street,
	}
}
