package app

import (
	"be2/internal/domain"
	"be2/internal/shared/date"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type fakeClientRepo struct {
	created   domain.Client
	deleted   uuid.UUID
	createErr error
	deleteErr error
}

func (f *fakeClientRepo) CreateClient(ctx context.Context, c domain.Client) error {
	f.created = c
	return f.createErr
}

func (f *fakeClientRepo) DeleteClient(ctx context.Context, id uuid.UUID) (int64, error) {
	f.deleted = id
	return 1, f.deleteErr
}

type fakeAddressService struct {
	createdInputs []CreateAddressInput
	idToReturn    uuid.UUID
	err           error
}

func (f *fakeAddressService) CreateAddress(ctx context.Context, c CreateAddressInput) (uuid.UUID, error) {
	f.createdInputs = append(f.createdInputs, c)
	return f.idToReturn, f.err
}

func (f *fakeAddressService) DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error) {
	return 0, nil
}
func (f *fakeAddressService) UpdateAddress(ctx context.Context, u UpdateInput) (int64, error) {
	return 0, nil
}

func TestCreateClientSuccess(t *testing.T) {
	addressID := uuid.New()
	addrSvc := &fakeAddressService{idToReturn: addressID}
	repo := &fakeClientRepo{}
	svc := NewClientService(repo, addrSvc)

	input := CreateClientInput{
		ClientName:       "Ivan",
		ClientSurname:    "Petrov",
		Birthday:         date.DateOnly{Time: time.Date(1990, 6, 1, 0, 0, 0, 0, time.UTC)},
		Gender:           domain.MALE,
		RegistrationDate: date.DateOnly{Time: time.Now().UTC()},
		Address:          CreateAddressInput{Country: "RU", City: "Moscow", Street: "Arbat"},
	}

	id, err := svc.CreateClient(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == uuid.Nil {
		t.Fatal("expected generated client ID")
	}
	if repo.created.Address != addressID {
		t.Fatalf("expected address %s, got %s", addressID, repo.created.Address)
	}
	if repo.created.ClientName != "Ivan" || repo.created.ClientSurname != "Petrov" {
		t.Fatalf("unexpected client data: %+v", repo.created)
	}
}

func TestCreateClientAddressError(t *testing.T) {
	addrSvc := &fakeAddressService{err: errors.New("addr error")}
	repo := &fakeClientRepo{}
	svc := NewClientService(repo, addrSvc)

	if _, err := svc.CreateClient(context.Background(), CreateClientInput{}); err == nil {
		t.Fatal("expected error from address service")
	}
}

func TestDeleteClient(t *testing.T) {
	repo := &fakeClientRepo{}
	svc := NewClientService(repo, &fakeAddressService{})
	clientID := uuid.New()

	deleted, err := svc.DeleteClient(context.Background(), clientID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleted != 1 || repo.deleted != clientID {
		t.Fatalf("unexpected delete result: %d, stored %s", deleted, repo.deleted)
	}
}
