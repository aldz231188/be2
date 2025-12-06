package app

import (
	"be2/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type fakeAddressRepo struct {
	lastCreated domainAddressRecord
	deletedID   uuid.UUID
	updated     domainAddressRecord
	createErr   error
	deleteErr   error
	updateErr   error
}

type domainAddressRecord struct {
	ID      uuid.UUID
	Country string
	City    string
	Street  string
}

func (f *fakeAddressRepo) CreateAddress(ctx context.Context, a domain.Address) error {
	f.lastCreated = domainAddressRecord{ID: a.ID, Country: a.Country, City: a.City, Street: a.Street}
	return f.createErr
}

func (f *fakeAddressRepo) DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error) {
	f.deletedID = id
	return 1, f.deleteErr
}

func (f *fakeAddressRepo) UpdateAddress(ctx context.Context, a domain.Address) (int64, error) {
	f.updated = domainAddressRecord{ID: a.ID, Country: a.Country, City: a.City, Street: a.Street}
	return 1, f.updateErr
}

func TestCreateAddressSuccess(t *testing.T) {
	repo := &fakeAddressRepo{}
	svc := NewAddressService(repo)

	id, err := svc.CreateAddress(context.Background(), CreateAddressInput{Country: "RU", City: "Moscow", Street: "Tverskaya"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == uuid.Nil {
		t.Fatal("expected generated UUID")
	}
	if repo.lastCreated.ID != id || repo.lastCreated.City != "Moscow" {
		t.Fatalf("address not saved correctly: %+v", repo.lastCreated)
	}
}

func TestCreateAddressContextError(t *testing.T) {
	repo := &fakeAddressRepo{createErr: context.Canceled}
	svc := NewAddressService(repo)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := svc.CreateAddress(ctx, CreateAddressInput{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}

func TestDeleteAddressErrorWrapped(t *testing.T) {
	repo := &fakeAddressRepo{deleteErr: errors.New("db error")}
	svc := NewAddressService(repo)

	if _, err := svc.DeleteAddress(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected error from delete")
	}
}

func TestUpdateAddress(t *testing.T) {
	repo := &fakeAddressRepo{}
	svc := NewAddressService(repo)
	id := uuid.New()

	updated, err := svc.UpdateAddress(context.Background(), UpdateInput{ID: id, Country: "RU", City: "Kazan", Street: "Kremlyovskaya"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated != 1 {
		t.Fatalf("expected 1 updated row, got %d", updated)
	}
	if repo.updated.ID != id || repo.updated.City != "Kazan" {
		t.Fatalf("unexpected updated record: %+v", repo.updated)
	}
}
