package app

import (
	"be2/services/bff/internal/domain"
	"be2/services/bff/internal/shared/date"

	"context"
	"errors"
	"fmt"

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
	DeleteClient(ctx context.Context, id uuid.UUID) (int64, error)
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
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return uuid.Nil, err
		}
		return uuid.Nil, fmt.Errorf("create client address: %w", err)
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
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return uuid.Nil, err
		}
		return uuid.Nil, fmt.Errorf("create client: %w", err)
	}

	return client.ID, nil
}

func (s *clientService) DeleteClient(ctx context.Context, id uuid.UUID) (int64, error) {
	deleted, err := s.clientRepo.DeleteClient(ctx, id)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return 0, err
		}
		return 0, fmt.Errorf("delete client: %w", err)
	}
	return deleted, nil
}
