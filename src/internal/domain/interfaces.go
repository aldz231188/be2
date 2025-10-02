package domain

import (
	"context"
)

// internal/domain/repo.go
type AddressRepo interface {
	CreateAddress(ctx context.Context, a Adress) error
}
type ClientRepo interface {
	CreateClient(ctx context.Context, c Client) error
	// GetClient(ctx context.Context, id uuid.UUID) (*Client, error)
}
