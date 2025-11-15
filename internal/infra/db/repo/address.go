package repo

import (
	"be2/internal/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
)

// type UserRepo struct{ q *Queries }

// type UpdateAddressParams struct {
// 	ID      uuid.UUID
// 	Country string
// 	City    string
// 	Street  string
// }

func (r *Repo) CreateAddress(ctx context.Context, c domain.Address) error {
	address := createAddressToRow(c)
	if err := r.q.CreateAddress(ctx, address); err != nil {
		return fmt.Errorf("create address: %w", err)
	}
	return nil
}
func (r *Repo) DeleteAddress(ctx context.Context, id uuid.UUID) (int64, error) {
	deleted, err := r.q.DeleteAddress(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("delete address: %w", err)
	}
	if deleted == 0 {
		return 0, domain.ErrAddressNotFound
	}
	return deleted, nil
}
func (r *Repo) UpdateAddress(ctx context.Context, c domain.Address) (int64, error) {
	address := updateAddressToRow(c)
	updated, err := r.q.UpdateAddress(ctx, address)
	if err != nil {
		return 0, fmt.Errorf("update address: %w", err)
	}
	if updated == 0 {
		return 0, domain.ErrAddressNotFound
	}
	return updated, nil
}
