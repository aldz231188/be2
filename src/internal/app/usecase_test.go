package app

import (
	"be2/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type repoStub struct {
	err      error
	received *domain.Adress
	calls    int
}

func (r *repoStub) AddAddress(ctx context.Context, addr *domain.Adress) error {
	r.calls++
	r.received = addr
	return r.err
}

func TestClientServiceImplAddAddressPropagatesRepoError(t *testing.T) {
	expectedErr := errors.New("db failure")
	repo := &repoStub{err: expectedErr}
	svc := NewClientServiceImpl(repo)

	err := svc.AddAddress(context.Background(), &domain.Adress{Country: "RU"})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if repo.calls != 1 {
		t.Fatalf("expected repo AddAddress to be called once, got %d", repo.calls)
	}
}

func TestClientServiceImplAddAddressSuccess(t *testing.T) {
	repo := &repoStub{}
	svc := NewClientServiceImpl(repo)

	input := &domain.Adress{Country: "RU", City: "Moscow", Street: "Tverskaya"}
	if err := svc.AddAddress(context.Background(), input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.calls != 1 {
		t.Fatalf("expected repo AddAddress to be called once, got %d", repo.calls)
	}

	if repo.received == nil {
		t.Fatal("expected repo to receive an address")
	}

	if repo.received == input {
		t.Fatal("expected service to pass a copied address instance")
	}

	if repo.received.Id == uuid.Nil {
		t.Error("expected generated address ID to be set")
	}

	if repo.received.Country != input.Country || repo.received.City != input.City || repo.received.Street != input.Street {
		t.Error("expected address fields to be copied to repo payload")
	}
}
