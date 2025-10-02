package domain

import (
	"context"

	"github.com/google/uuid"
)

// internal/domain/repo.go
type AddressRepo interface {
	CreateAddress(ctx context.Context, a Address) error
	DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error)
}
type ClientRepo interface {
	CreateClient(ctx context.Context, c Client) error
	// GetClient(ctx context.Context, id uuid.UUID) (*Client, error)
}
