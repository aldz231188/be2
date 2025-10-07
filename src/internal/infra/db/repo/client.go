package repo

import (
	"be2/internal/domain"
	"context"
	// "github.com/google/uuid"
)

// type UserRepo struct{ q *Queries }

// type UpdateAddressParams struct {
// 	ID      uuid.UUID
// 	Country interface{}
// 	City    interface{}
// 	Street  interface{}
// }

func (r *Repo) CreateClient(ctx context.Context, c domain.Client) error {
	client := createClientToRow(c)
	return r.q.CreateClient(ctx, client)
}

// func (r *Repo) DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error) {
// 	return r.q.DeleteAddress(ctx, id)
// }
// func (r *Repo) UpdateAddress(ctx context.Context, c domain.Address) (int64, error) {
// 	address := updateAddressToRow(c)
// 	return r.q.UpdateAddress(ctx, address)
// }

// реализация CustomerRepo
// func (r *Repo) CreateCustomer(ctx context.Context, c domain.Client) error             { /* ... */ }
// func (r *Repo) GetCustomer(ctx context.Context, id uuid.UUID) (*domain.Client, error) { /* ... */ }

// var _ domain.CustomerRepo = (*Repo)(nil)
// var _ domain.SupplierRepo = (*Repo)(nil)
