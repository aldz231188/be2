package db

import (
	"be2/internal/domain"
	// "be2/internal/infra/db/postgres"
	"context"
)

// type UserRepo struct{ q *Queries }

func (r *UserRepo) AddAddress(ctx context.Context, c *domain.Adress) error {

	address := domaineToRow(c)
	r.q.AddAddress(ctx, *address)
	return nil

}
