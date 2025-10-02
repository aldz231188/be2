package app

import (
	"be2/internal/domain"

	"context"
	// "errors"
	"github.com/google/uuid"
	// "time"
)

type CreateInput struct {
	Country string
	City    string
	Street  string
}

type Service interface {
	CreateAddress(ctx context.Context, c CreateInput) (uuid.UUID, error)
	DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error)
}

var _ Service = (*ServiceImpl)(nil) // гарантирует, что *ServiceImpl реализует Service

type ServiceImpl struct { //вынести
	repo domain.AddressRepo
}

func NewServiceImpl(r domain.AddressRepo) Service { //вынести
	return &ServiceImpl{repo: r}
}

func (si *ServiceImpl) CreateAddress(ctx context.Context, c CreateInput) (uuid.UUID, error) {
	address := domain.Address{
		Id:      uuid.New(),
		Country: c.Country,
		City:    c.City,
		Street:  c.Street,
	}

	if err := si.repo.CreateAddress(ctx, address); err != nil {
		return uuid.Nil, err
	}

	return address.Id, nil

}

func (si *ServiceImpl) DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error) {
	return si.repo.DeleteAddress(ctx, id)
}

// func NormalizeDate(d time.Time) time.Time {
// 	y, m, day := d.Date()
// 	return time.Date(y, m, day, 0, 0, 0, 0, time.UTC)
// }
