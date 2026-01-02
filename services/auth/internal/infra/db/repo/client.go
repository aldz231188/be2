package repo

import (
	"be2/services/auth/internal/domain"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repo) CreateClient(ctx context.Context, c domain.Client) error {
	client := createClientToRow(c)
	if err := r.q.CreateClient(ctx, client); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrClientAlreadyExists
		}
		return fmt.Errorf("create client: %w", err)
	}
	return nil
}

func (r *Repo) DeleteClient(ctx context.Context, id uuid.UUID) (int64, error) {
	deleted, err := r.q.DeleteClient(ctx, id)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return 0, err
		}
		return 0, fmt.Errorf("delete client: %w", err)
	}
	if deleted == 0 {
		return 0, domain.ErrClientNotFound
	}
	return deleted, nil
}
