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
