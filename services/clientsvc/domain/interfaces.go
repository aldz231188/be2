package domain

import (
	"context"

	"github.com/google/uuid"
)

// internal/domain/repo.go
//
//	type AddressRepo interface {
//		CreateAddress(ctx context.Context, a Address) error
//		DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error)
//		UpdateAddress(ctx context.Context, a Address) (int64, error)
//	}
type ClientRepo interface {
	CreateClient(ctx context.Context, a Client) error
	DeleteClient(ctx context.Context, id uuid.UUID) (int64, error)
	// UpdateClient(ctx context.Context, a Address) (int64, error)
}
