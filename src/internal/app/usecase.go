package app

import (
	"be2/internal/domain"

	"context"
	// "errors"
	"github.com/google/uuid"
	"time"
)

type ClientServiceImpl struct {
	repo domain.Repo
}

func NewClientServiceImpl(r domain.Repo) domain.Service {
	return &ClientServiceImpl{repo: r}
}

func (csi *ClientServiceImpl) AddAddress(ctx context.Context, c *domain.Adress) error {
	adress, err := NewAdress(c)
	if err != nil {
		return err
	}

	// birthday := NormalizeDate(c.Birthday)
	// if birthday.After(NormalizeDate(time.Now().UTC())) {
	// 	return nil, errors.New("birthday in the future")
	// }

	// client := domain.Client{
	// 	Id:               uuid.New(),
	// 	ClientName:       c.ClientName,
	// 	ClientSurname:    c.ClientSurname,
	// 	Birthday:         birthday,
	// 	Gender:           c.Gender,
	// 	RegistrationDate: time.Now(),
	// 	Address:          adress,
	// }
	csi.repo.AddAddress(ctx, adress)

	return nil

}

func NewAdress(a *domain.Adress) (*domain.Adress, error) {
	return &domain.Adress{
		Id:      uuid.New(),
		Country: a.Country,
		City:    a.City,
		Street:  a.Street,
	}, nil

}

func NormalizeDate(d time.Time) time.Time {
	y, m, day := d.Date()
	return time.Date(y, m, day, 0, 0, 0, 0, time.UTC)
}
