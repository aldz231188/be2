package app

import (
	"be2/services/auth/internal/domain"

	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
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

type addressService struct {
	addressRepo domain.AddressRepo
}

var _ AddressService = (*addressService)(nil)

func NewAddressService(repo domain.AddressRepo) AddressService {
	return &addressService{addressRepo: repo}
}

func (s *addressService) CreateAddress(ctx context.Context, c CreateAddressInput) (uuid.UUID, error) {
	address := domain.Address{
		ID:      uuid.New(),
		Country: c.Country,
		City:    c.City,
		Street:  c.Street,
	}

	if err := s.addressRepo.CreateAddress(ctx, address); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return uuid.Nil, err
		}
		return uuid.Nil, fmt.Errorf("create address: %w", err)
	}

	return address.ID, nil
}

func (s *addressService) DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error) {
	deleted, err := s.addressRepo.DeleteAddress(ctx, id)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return 0, err
		}
		return 0, fmt.Errorf("delete address: %w", err)
	}
	return deleted, nil
}

func (s *addressService) UpdateAddress(ctx context.Context, u UpdateInput) (int64, error) {
	address := domain.Address{
		ID:      u.ID,
		Country: u.Country,
		City:    u.City,
		Street:  u.Street,
	}
	updated, err := s.addressRepo.UpdateAddress(ctx, address)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return 0, err
		}
		return 0, fmt.Errorf("update address: %w", err)
	}
	return updated, nil
}
