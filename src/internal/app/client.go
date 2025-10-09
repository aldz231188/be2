package app

import (
	"be2/internal/domain"
	"be2/internal/shared/date"

	"context"

	"github.com/google/uuid"
)

type CreateClientInput struct {
	ClientName       string
	ClientSurname    string
	Birthday         date.DateOnly
	Gender           domain.Gender
	RegistrationDate date.DateOnly
	Address          CreateAddressInput
}

type ClientService interface {
	CreateClient(ctx context.Context, c CreateClientInput) (uuid.UUID, error)
}

type clientService struct {
	clientRepo     domain.ClientRepo
	addressService AddressService
}

var _ ClientService = (*clientService)(nil)

func NewClientService(clientRepo domain.ClientRepo, addressSvc AddressService) ClientService {
	return &clientService{
		clientRepo:     clientRepo,
		addressService: addressSvc,
	}
}

func (s *clientService) CreateClient(ctx context.Context, c CreateClientInput) (uuid.UUID, error) {
	addressID, err := s.addressService.CreateAddress(ctx, c.Address)
	if err != nil {
		return uuid.Nil, err
	}

	client := domain.Client{
		ID:               uuid.New(),
		ClientName:       c.ClientName,
		ClientSurname:    c.ClientSurname,
		Birthday:         c.Birthday.Time,
		Gender:           c.Gender,
		RegistrationDate: c.RegistrationDate.Time,
		Address:          addressID,
	}

	if err := s.clientRepo.CreateClient(ctx, client); err != nil {
		return uuid.Nil, err
	}

	return client.ID, nil
}
