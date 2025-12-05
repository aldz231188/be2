package domain

import (
	"context"

	"github.com/google/uuid"
)

// internal/domain/repo.go
type AddressRepo interface {
	CreateAddress(ctx context.Context, a Address) error
	DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error)
	UpdateAddress(ctx context.Context, a Address) (int64, error)
}
type ClientRepo interface {
	CreateClient(ctx context.Context, a Client) error
	DeleteClient(ctx context.Context, id uuid.UUID) (int64, error)
	// UpdateClient(ctx context.Context, a Address) (int64, error)
}

type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	IncrementTokenVersion(ctx context.Context, id uuid.UUID) error
}

type SessionRepo interface {
	CreateSession(ctx context.Context, session Session) error
	GetSessionByHash(ctx context.Context, hash string) (Session, error)
	RevokeSession(ctx context.Context, hash string) error
	RevokeSessionsByUser(ctx context.Context, userID uuid.UUID) error
}
