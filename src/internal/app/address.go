package app

import (
	"be2/internal/domain"

	"context"
	// "errors"
	"github.com/google/uuid"
	// "time"
)

type CreateAddressInput struct {
	Country string
	City    string
	Street  string
}
type UpdateInput struct {
	ID      uuid.UUID
	Country string
	City    string
	Street  string
}

type AddressService interface {
	CreateAddress(ctx context.Context, c CreateAddressInput) (uuid.UUID, error)
	DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error)
	UpdateAddress(ctx context.Context, u UpdateInput) (int64, error)
}

var _ AddressService = (*ServiceImpl)(nil) // гарантирует, что *ServiceImpl реализует AddressService
var _ ClientService = (*ServiceImpl)(nil)  // гарантирует, что *ServiceImpl реализует ClientService

type ServiceImpl struct { //вынести
	addressrepo domain.AddressRepo
	clientrepo  domain.ClientRepo
}

func NewServiceImpl(a domain.AddressRepo, c domain.ClientRepo) *ServiceImpl { //вынести
	return &ServiceImpl{
		addressrepo: a,
		clientrepo:  c,
	}
}

func (si *ServiceImpl) CreateAddress(ctx context.Context, c CreateAddressInput) (uuid.UUID, error) {
	address := domain.Address{
		ID:      uuid.New(),
		Country: c.Country,
		City:    c.City,
		Street:  c.Street,
	}

	if err := si.addressrepo.CreateAddress(ctx, address); err != nil {
		return uuid.Nil, err
	}

	return address.ID, nil

}

func (si *ServiceImpl) DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error) {
	return si.addressrepo.DeleteAddress(ctx, id)
}
func (si *ServiceImpl) UpdateAddress(ctx context.Context, u UpdateInput) (int64, error) {
	address := domain.Address{
		ID:      u.ID,
		Country: u.Country,
		City:    u.City,
		Street:  u.Street,
	}
	return si.addressrepo.UpdateAddress(ctx, address)
}

// func NormalizeDate(d time.Time) time.Time {
// 	y, m, day := d.Date()
// 	return time.Date(y, m, day, 0, 0, 0, 0, time.UTC)
// }
