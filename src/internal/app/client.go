package app

import (
	"be2/internal/domain"
	"be2/internal/shared/date"

	"context"
	// "errors"
	"github.com/google/uuid"
	// "time"
)

type CreateClientInput struct {
	ClientName       string
	ClientSurname    string
	Birthday         date.DateOnly
	Gender           domain.Gender
	RegistrationDate date.DateOnly
	Address          CreateAddressInput
}

// type UpdateInput struct {
// 	ID      uuid.UUID
// 	Country string
// 	City    string
// 	Street  string
// }

type ClientService interface {
	CreateClient(ctx context.Context, c CreateClientInput) (uuid.UUID, error)
	// DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error)
	// UpdateAddress(ctx context.Context, u UpdateInput) (int64, error)
}

// var _ ClientService = (*ServiceImpl)(nil) // гарантирует, что *ServiceImpl реализует Service

// type ServiceImpl struct { //вынести
// 	repo domain.AddressRepo
// }

// func NewServiceImpl(r domain.AddressRepo) Service { //вынести
// 	return &ServiceImpl{repo: r}
// }

func (si *ServiceImpl) CreateClient(ctx context.Context, c CreateClientInput) (uuid.UUID, error) {

	if address, err := si.CreateAddress(ctx, c.Address); err != nil {
		return uuid.Nil, err
	} else {
		client := domain.Client{
			ID:            uuid.New(),
			ClientName:    c.ClientName,
			ClientSurname: c.ClientSurname,
			Birthday:      c.Birthday.Time,
			Gender:        c.Gender,
			Address:       address,
		}
		si.clientrepo.CreateClient(ctx, client)
		return client.ID, nil
	}

}

// func (si *ServiceImpl) DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error) {
// 	return si.repo.DeleteAddress(ctx, id)
// }
// func (si *ServiceImpl) UpdateAddress(ctx context.Context, u UpdateInput) (int64, error) {
// 	address := domain.Address{
// 		ID:      u.ID,
// 		Country: u.Country,
// 		City:    u.City,
// 		Street:  u.Street,
// 	}
// 	return si.repo.UpdateAddress(ctx, address)
// }

// // func NormalizeDate(d time.Time) time.Time {
// // 	y, m, day := d.Date()
// // 	return time.Date(y, m, day, 0, 0, 0, 0, time.UTC)
// // }
