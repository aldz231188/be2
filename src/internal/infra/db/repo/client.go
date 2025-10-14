package repo

import (
	"be2/internal/domain"
	"context"
	"github.com/google/uuid"
)

func (r *Repo) CreateClient(ctx context.Context, c domain.Client) error {
	client := createClientToRow(c)
	return r.q.CreateClient(ctx, client)
}

func (r *Repo) DeleteClient(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.q.DeleteClient(ctx, id)
}

// func (r *Repo) UpdateAddress(ctx context.Context, c domain.Address) (int64, error) {
// 	address := updateAddressToRow(c)
// 	return r.q.UpdateAddress(ctx, address)
// }

// реализация CustomerRepo
// func (r *Repo) CreateCustomer(ctx context.Context, c domain.Client) error             { /* ... */ }
// func (r *Repo) GetCustomer(ctx context.Context, id uuid.UUID) (*domain.Client, error) { /* ... */ }

// var _ domain.CustomerRepo = (*Repo)(nil)
// var _ domain.SupplierRepo = (*Repo)(nil)
