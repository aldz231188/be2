package db

import (
	"be2/internal/domain"
	// "be2/internal/infra/db/postgres"
	"context"
)

// type UserRepo struct{ q *Queries }

func (r *Repo) CreateAddress(ctx context.Context, c domain.Adress) error {
	address := domaineToRow(c)
	return r.q.AddAddress(ctx, address)
}

// реализация CustomerRepo
// func (r *Repo) CreateCustomer(ctx context.Context, c domain.Client) error             { /* ... */ }
// func (r *Repo) GetCustomer(ctx context.Context, id uuid.UUID) (*domain.Client, error) { /* ... */ }

// var _ domain.CustomerRepo = (*Repo)(nil)
// var _ domain.SupplierRepo = (*Repo)(nil)
