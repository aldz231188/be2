package repo

import (
	"be2/internal/domain"
	"context"
	"github.com/google/uuid"
)

// type UserRepo struct{ q *Queries }

func (r *Repo) CreateAddress(ctx context.Context, c domain.Address) error {
	address := domaineToRow(c)
	return r.q.CreateAddress(ctx, address)
}
func (r *Repo) DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.q.DeleteAddress(ctx, id)
}

// реализация CustomerRepo
// func (r *Repo) CreateCustomer(ctx context.Context, c domain.Client) error             { /* ... */ }
// func (r *Repo) GetCustomer(ctx context.Context, id uuid.UUID) (*domain.Client, error) { /* ... */ }

// var _ domain.CustomerRepo = (*Repo)(nil)
// var _ domain.SupplierRepo = (*Repo)(nil)
